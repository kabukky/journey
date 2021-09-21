//go:build !noplugins
// +build !noplugins

package plugins

import (
	"github.com/Landria/journey/structure"
	"github.com/Landria/journey/structure/methods"
	lua "github.com/yuin/gopher-lua"
)

func convertArguments(vm *lua.LState, structureArguments []structure.Helper) *lua.LTable {
	table := vm.NewTable()
	arguments := methods.ProcessHelperArguments(structureArguments)
	for key, value := range arguments {
		table.RawSet(lua.LString(key), lua.LString(value))
	}
	return table
}

func convertPost(vm *lua.LState, structurePost *structure.Post) *lua.LTable {
	post := vm.NewTable()
	post.RawSet(lua.LString("id"), lua.LNumber(structurePost.Id))
	post.RawSet(lua.LString("uuid"), lua.LString(structurePost.Uuid))
	post.RawSet(lua.LString("title"), lua.LString(structurePost.Title))
	post.RawSet(lua.LString("slug"), lua.LString(structurePost.Slug))
	post.RawSet(lua.LString("markdown"), lua.LString(structurePost.Markdown))
	post.RawSet(lua.LString("html"), lua.LString(structurePost.Html))
	post.RawSet(lua.LString("isfeatured"), lua.LBool(structurePost.IsFeatured))
	post.RawSet(lua.LString("ispage"), lua.LBool(structurePost.IsPage))
	post.RawSet(lua.LString("ispublished"), lua.LBool(structurePost.IsPublished))
	post.RawSet(lua.LString("date"), lua.LNumber(structurePost.Date.Unix()))
	post.RawSet(lua.LString("image"), lua.LString(structurePost.Image))
	post.RawSet(lua.LString("metadescription"), lua.LString(structurePost.MetaDescription))
	return post
}

func convertUser(vm *lua.LState, structureUser *structure.User) *lua.LTable {
	user := vm.NewTable()
	user.RawSet(lua.LString("id"), lua.LNumber(structureUser.Id))
	user.RawSet(lua.LString("name"), lua.LString(structureUser.Name))
	user.RawSet(lua.LString("slug"), lua.LString(structureUser.Slug))
	user.RawSet(lua.LString("email"), lua.LString(structureUser.Email))
	user.RawSet(lua.LString("image"), lua.LString(structureUser.Image))
	user.RawSet(lua.LString("cover"), lua.LString(structureUser.Cover))
	user.RawSet(lua.LString("bio"), lua.LString(structureUser.Bio))
	user.RawSet(lua.LString("website"), lua.LString(structureUser.Website))
	user.RawSet(lua.LString("location"), lua.LString(structureUser.Location))
	user.RawSet(lua.LString("role"), lua.LNumber(structureUser.Role))
	return user
}

func convertTags(vm *lua.LState, structureTags []structure.Tag) *lua.LTable {
	table := make([]*lua.LTable, 0)
	for index, _ := range structureTags {
		tag := vm.NewTable()
		tag.RawSet(lua.LString("id"), lua.LNumber(structureTags[index].Id))
		tag.RawSet(lua.LString("name"), lua.LString(structureTags[index].Name))
		tag.RawSet(lua.LString("slug"), lua.LString(structureTags[index].Slug))
		table = append(table, tag)
	}
	return makeTable(vm, table)
}

func convertBlog(vm *lua.LState, structureBlog *structure.Blog) *lua.LTable {
	blog := vm.NewTable()
	blog.RawSet(lua.LString("url"), lua.LString(structureBlog.Url))
	blog.RawSet(lua.LString("title"), lua.LString(structureBlog.Title))
	blog.RawSet(lua.LString("description"), lua.LString(structureBlog.Description))
	blog.RawSet(lua.LString("logo"), lua.LString(structureBlog.Logo))
	blog.RawSet(lua.LString("cover"), lua.LString(structureBlog.Cover))
	blog.RawSet(lua.LString("assetpath"), lua.LString(structureBlog.AssetPath))
	blog.RawSet(lua.LString("postcount"), lua.LNumber(structureBlog.PostCount))
	blog.RawSet(lua.LString("postsperpage"), lua.LNumber(structureBlog.PostsPerPage))
	blog.RawSet(lua.LString("activetheme"), lua.LString(structureBlog.ActiveTheme))
	return blog
}

func makeTable(vm *lua.LState, tables []*lua.LTable) *lua.LTable {
	table := vm.NewTable()
	for index, _ := range tables {
		table.Append(tables[index])
	}
	return table
}
