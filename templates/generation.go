package templates

import (
	"bytes"
	"errors"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/structure"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var openTag = []byte("{{")
var closeTag = []byte("}}")

func getFunction(name string) func(*Helper, *structure.RequestData) []byte {
	if helperFuctions[name] != nil {
		return helperFuctions[name]
	} else {
		return helperFuctions["null"]
	}
}

func createHelper(helperName []byte, unescaped bool, startPos int, block []byte, children []Helper, elseHelper *Helper) *Helper {
	var helper *Helper
	// Check for =arguments
	twoPartArgumentChecker := regexp.MustCompile("(\\S+?)\\s*?=\\s*?['\"](.*?)['\"]")
	twoPartArgumentResult := twoPartArgumentChecker.FindAllSubmatch(helperName, -1)
	twoPartArguments := make([][]byte, 0)
	for _, arg := range twoPartArgumentResult {
		if len(arg) == 3 {
			twoPartArguments = append(twoPartArguments, bytes.Join(arg[1:], []byte("")))
			//remove =argument from helper name
			helperName = bytes.Replace(helperName, arg[0], []byte(""), 1)
		}
	}
	// Separate arguments (e.g. 'if @blog.title')
	tags := bytes.Fields(helperName)
	for index, tag := range tags {
		//remove "" around tag if present
		quoteTagChecker := regexp.MustCompile("^[\"'](.+?)[\"']$")
		quoteTagResult := quoteTagChecker.FindSubmatch(tag)
		if len(quoteTagResult) != 0 {
			tag = quoteTagResult[1]
		}
		//TODO: This may have to change if the first argument is surrounded by ""
		if index == 0 {
			helper = makeHelper(string(tag), unescaped, startPos, block, children)
		} else {
			// Handle whitespaces in arguments
			helper.Arguments = append(helper.Arguments, *makeHelper(string(tag), unescaped, 0, []byte{}, nil))
		}
	}
	if len(twoPartArguments) != 0 {
		for _, arg := range twoPartArguments {
			helper.Arguments = append(helper.Arguments, *makeHelper(string(arg), unescaped, 0, []byte{}, nil))
		}
	}
	if elseHelper != nil {
		helper.Arguments = append(helper.Arguments, *elseHelper)
	}
	return helper
}

func makeHelper(tag string, unescaped bool, startPos int, block []byte, children []Helper) *Helper {
	return &Helper{Name: tag, Arguments: nil, Unescaped: unescaped, Position: startPos, Block: block, Children: children, Function: getFunction(tag)}
}

func findHelper(data []byte, allHelpers []Helper) ([]byte, []Helper) {
	startPos := bytes.Index(data, openTag)
	endPos := bytes.Index(data, closeTag)
	if startPos != -1 && endPos != -1 {
		openTagLength := len(openTag)
		closeTagLength := len(closeTag)
		unescaped := false
		helperName := data[startPos+openTagLength : endPos]
		// Check if helper calls for unescaped text (e.g. three brackets - {{{title}}})
		if bytes.HasPrefix(helperName, []byte("{")) {
			unescaped = true
			openTagLength++ //not necessary
			closeTagLength++
			helperName = helperName[len([]byte("{")):]
		}
		helperName = bytes.Trim(helperName, " ") //make sure there are no trailing whitespaces
		// Remove helper from data
		parts := [][]byte{data[:startPos], data[endPos+closeTagLength:]}
		data = bytes.Join(parts, []byte(""))
		// Check if comment
		if bytes.HasPrefix(helperName, []byte("! ")) || bytes.HasPrefix(helperName, []byte("!--")) {
			return findHelper(data, allHelpers)
		}
		// Check if block
		if bytes.HasPrefix(helperName, []byte("#")) {
			helperName = helperName[len([]byte("#")):] //remove '#' from helperName
			var helper Helper
			data, helper = findBlock(data, helperName, unescaped, startPos) //only use the data string after the opening tag
			allHelpers = append(allHelpers, helper)
			return findHelper(data, allHelpers)
		}
		allHelpers = append(allHelpers, *createHelper(helperName, unescaped, startPos, []byte{}, nil, nil))
		return findHelper(data, allHelpers)
	} else {
		return data, allHelpers
	}
}

func findBlock(data []byte, helperName []byte, unescaped bool, startPos int) ([]byte, Helper) {
	arguments := bytes.Fields(helperName)
	tag := arguments[0] // Get only the first tag (e.g. 'if' in 'if @blog.cover')
	arguments = arguments[1:]
	closeParts := []string{"{{2,3}\\s*/", string(tag), ".?}{2,3}"}
	openParts := []string{"{{2,3}\\s*#", string(tag), ".+?}{2,3}"}
	closeRegex := regexp.MustCompile(strings.Join(closeParts, ""))
	openRegex := regexp.MustCompile(strings.Join(openParts, ""))
	closePositions := closeRegex.FindAllIndex(data, -1)
	openPositions := openRegex.FindAllIndex(data, -1)
	// Check if there are opening tags before the closing tag
	positionIndex := 0
	for _, openPosition := range openPositions {
		if openPosition[0] < closePositions[positionIndex][0] {
			positionIndex++
		}
	}
	block := data[startPos:closePositions[positionIndex][0]]
	parts := [][]byte{data[:startPos], data[closePositions[positionIndex][1]:]}
	data = bytes.Join(parts, []byte(""))
	children := make([]Helper, 0)
	block, children = findHelper(block, children)
	// Handle else (search children for else helper)
	for index, child := range children {
		if child.Name == "else" {
			elseHelper := child
			// Change blocks
			elseHelper.Block = block[elseHelper.Position:]
			block = block[:elseHelper.Position]
			// Change children, omit else helper
			elseHelper.Children = children[(index + 1):]
			// Change Position in children of else helper
			for indexElse, _ := range elseHelper.Children {
				elseHelper.Children[indexElse].Position = elseHelper.Children[indexElse].Position - elseHelper.Position
			}
			children = children[:index]
			helper := createHelper(helperName, unescaped, startPos, block, children, &elseHelper)
			return data, *helper
		}
	}
	helper := createHelper(helperName, unescaped, startPos, block, children, nil)
	return data, *helper
}

func compileTemplate(data []byte, name string) *Helper {
	baseHelper := Helper{Name: name, Arguments: nil, Unescaped: false, Position: 0, Block: []byte{}, Children: nil, Function: getFunction(name)}
	allHelpers := make([]Helper, 0)
	data, allHelpers = findHelper(data, allHelpers)
	baseHelper.Block = data
	baseHelper.Children = allHelpers
	// Handle extend helper
	for index, child := range baseHelper.Children {
		if child.Name == "body" {
			baseHelper.BodyHelper = &baseHelper.Children[index]
		}
	}
	return &baseHelper
}

func createTemplateFromFile(filename string) (*Helper, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	helper := compileTemplate(data, path.Base(filename)[0:len(path.Base(filename))-len(path.Ext(filename))]) //second argument: get filename without extension
	return helper, nil
}

func inspectTemplateFile(filePath string, info os.FileInfo, err error) error {
	if !info.IsDir() && path.Ext(filePath) == ".hbs" {
		helper, err := createTemplateFromFile(filePath)
		if err != nil {
			return err
		}
		compiledTemplates.m[helper.Name] = helper
	}
	return nil
}

func Generate() error {
	compiledTemplates.Lock()
	defer compiledTemplates.Unlock()
	activeTheme, err := database.RetrieveActiveTheme()
	if err != nil {
		return err
	}
	// Compile all template files
	// First clear compiledTemplates map (theme could have been changed) TODO: Should this be implemented?
	currentThemePath := path.Join(filenames.ThemesFilepath, *activeTheme)
	err = filepath.Walk(currentThemePath, inspectTemplateFile)
	if err != nil {
		return err
	}
	// Check if index and post templates are compiled
	if _, ok := compiledTemplates.m["index"]; !ok {
		return errors.New("Couldn't compile template 'index'. Is index.hbs missing?")
	}
	if _, ok := compiledTemplates.m["post"]; !ok {
		return errors.New("Couldn't compile template 'post'. Is post.hbs missing?")
	}
	return nil
}

func GetAllThemes() []string {
	themes := make([]string, 0)
	files, _ := filepath.Glob(filepath.Join(filenames.ThemesFilepath, "*"))
	for _, file := range files {
		if isDirectory(file) {
			themes = append(themes, filepath.Base(file))
		}
	}
	return themes
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
