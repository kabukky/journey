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
	// Make map
	nameMap := make(map[string]string, 0)
	err := filepath.Walk(filenames.PluginsFilepath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(filePath) == ".lua" {
			// Check if the lua file is a plugin entry point by executing it
			helperNames, err := getHelperNames(filePath)
			if err != nil {
				return err
			}
			// Add all file names of helpers to the name map
			for _, helperName := range helperNames {
				log.Println("Helper name:", helperName)
				absPath, err := filepath.Abs(filePath)
				if err != nil {
					log.Println("Error while determining absolute path to lua file:", err)
					return err
				}
				nameMap[helperName] = absPath
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(nameMap) == 0 {
		return errors.New("No plugins were loaded.")
	}
	// If plugins were loaded, create LuaPool and assign name map to LuaPool
	LuaPool = newLuaPool()
	LuaPool.m.Lock()
	defer LuaPool.m.Unlock()
	LuaPool.files = nameMap
	return nil
}

func getHelperNames(fileName string) ([]string, error) {
	// Make a slice to hold all helper names
	helperList := make([]string, 0)
	// Create a new lua state
	vm := lua.NewState()
	defer vm.Close()
	// Set up vm functions
	values := &structure.RequestData{}
	absDir, err := filepath.Abs(fileName)
	if err != nil {
		log.Println("Error while determining absolute path to lua file:", err)
		return helperList, err
	}
	setUpVm(vm, values, absDir)
	// Execute plugin
	// TODO: Is there a better way to just load the file? We only need to execute the register function (see below)
	err = vm.DoFile(absDir)
	if err != nil {
		// TODO: We are not returning upon error here. Keep it like this?
		log.Println("Error while loading plugin:", err)
	}
	err = vm.CallByParam(lua.P{Fn: vm.GetGlobal("register"), NRet: 1, Protect: true})
	if err != nil {
		// Fail silently since this is probably just a lua file without a register function
		return helperList, nil
	}
	// Get return value
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
	}
	return helperList, nil
}

// Creates all methods that can be used from Lua.
func setUpVm(vm *lua.LState, values *structure.RequestData, absPathToLuaFile string) {
	luaPath := filepath.Dir(absPathToLuaFile)
	// Function to get the directory of the current file (to add to LUA_PATH in Lua)
	vm.SetGlobal("getCurrentDir", vm.NewFunction(func(vm *lua.LState) int {
		vm.Push(lua.LString(luaPath))
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
