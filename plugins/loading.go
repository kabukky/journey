package plugins

import (
	"errors"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/structure"
	"github.com/yuin/gopher-lua"
	"log"
	"os"
	"path/filepath"
)

func Load() error {
	// Clear LuaPool
	LuaPool = nil
	// Make map
	stateMap := make(map[string]*lua.LState, 0)
	err := filepath.Walk(filenames.PluginsFilepath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(filePath) == ".lua" {
			// Initialize the LuaPool if it is not alread initialized
			if LuaPool == nil {
				makeLuaPool()
			}
			// Lock LuaPool
			LuaPool.Lock()
			defer LuaPool.Unlock()
			// Check if the lua file is a plugin entry point by executing it
			helperNames, vm := getHelperNames(filePath)
			for _, helperName := range helperNames {
				log.Println("Helper name:", helperName)
				stateMap[helperName] = vm
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(stateMap) == 0 {
		LuaPool = nil
		return errors.New("No plugins were loaded.")
	}
	// If plugins were loaded, assign state map to the first element of the array in LuaPool
	LuaPool.Lock()
	defer LuaPool.Unlock()
	LuaPool.template = stateMap
	return nil
}

func getHelperNames(fileName string) ([]string, *lua.LState) {
	// Make a slice to hold all helper names
	helperList := make([]string, 0)
	// Create a new lua state
	vm := lua.NewState()
	// Set up vm (make sure to append the absolute path to the Lua script to LUA_PATH)
	values := &structure.RequestData{}
	absDir, err := filepath.Abs(fileName)
	setUpVm(vm, values, filepath.Dir(absDir))
	// Execute plugin
	err = vm.DoFile(fileName)
	if err != nil {
		log.Println("Error while loading plugin:", err)
	}
	err = vm.CallByParam(lua.P{Fn: vm.GetGlobal("register"), NRet: 1, Protect: true})
	if err != nil {
		vm.Close()
		return helperList, nil
	}
	table := vm.ToTable(-1)
	// Check if return value is a table
	if table != nil {
		// Iterate the table for every helper name to be registered
		table.ForEach(func(key lua.LValue, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				if string(str) != "" {
					helperList = append(helperList, string(str))
				}
			}
		})
	} else { // Else return nil
		vm.Close()
		return []string{}, nil
	}
	return helperList, vm
}

// Creates all methods that can be used by lua. The isTesting argument indicates if it should be set up just for testing (e. g. loading and initial testing of the lua file).
func setUpVm(vm *lua.LState, values *structure.RequestData, filePath string) {
	// Function to get the dir of the current file (to add to LUA_PATH in Lua)
	vm.SetGlobal("getCurrentDir", vm.NewFunction(func(vm *lua.LState) int {
		vm.Push(lua.LString(filePath))
		return 1 // Number of results
	}))
	// Function to print to the log
	vm.SetGlobal("print", vm.NewFunction(func(vm *lua.LState) int {
		log.Println(vm.Get(-1).String())
		return 0 // Number of results
	}))
	// Function to get number of posts in values
	vm.SetGlobal("getNumberOfPosts", vm.NewFunction(func(vm *lua.LState) int {
		vm.Push(lua.LNumber(len(values.Posts)))
		return 1 // Number of results
	}))
	// Function to get a post by its index
	vm.SetGlobal("getPost", vm.NewFunction(func(vm *lua.LState) int {
		postIndex := vm.ToInt(-1)
		vm.Push(convertPost(vm, &values.Posts[postIndex-1]))
		return 1 // Number of results
	}))
	// Function to get a user by post
	vm.SetGlobal("getAuthorForPost", vm.NewFunction(func(vm *lua.LState) int {
		postIndex := vm.ToInt(-1)
		vm.Push(convertUser(vm, values.Posts[postIndex-1].Author))
		return 1 // Number of results
	}))
	// Function to get tags by post
	vm.SetGlobal("getTagsForPost", vm.NewFunction(func(vm *lua.LState) int {
		postIndex := vm.ToInt(-1)
		vm.Push(convertTags(vm, values.Posts[postIndex-1].Tags))
		return 1 // Number of results
	}))
	// Function to get blog
	vm.SetGlobal("getBlog", vm.NewFunction(func(vm *lua.LState) int {
		vm.Push(convertBlog(vm, values.Blog))
		return 1 // Number of results
	}))
}
