package templates

import (
	"github.com/kabukky/journey/structure"
)

var helperFuctions = map[string]func(*structure.Helper, *structure.RequestData) []byte{

	// Null function
	"null": nullFunc,

	// General functions
	"if":               ifFunc,
	"unless":           unlessFunc,
	"foreach":          foreachFunc,
	"!<":               extendFunc,
	"body":             bodyFunc,
	"asset":            assetFunc,
	"encode":           encodeFunc,
	">":                insertFunc,
	"meta_title":       meta_titleFunc,
	"meta_description": meta_descriptionFunc,
	"ghost_head":       ghost_headFunc,
	"ghost_foot":       ghost_footFunc,
	"body_class":       body_classFunc,
	"plural":           pluralFunc,
	"date":             dateFunc,
	"image":            imageFunc,
	"contentFor":       contentForFunc,
	"block":            blockFunc,

	// TODO: remove legacy API v2
	// @blog functions
	"@blog.title":       atBlogDotTitleFunc,
	"@blog.url":         atBlogDotUrlFunc,
	"@blog.logo":        atBlogDotLogoFunc,
	"@blog.cover":       atBlogDotCoverFunc,
	"@blog.cover_image": atBlogDotCoverFunc,
	"@blog.description": atBlogDotDescriptionFunc,
	"@blog.navigation":  navigationFunc,

	// @site functions
	"@site.title":       atBlogDotTitleFunc,
	"@site.url":         atBlogDotUrlFunc,
	"@site.logo":        atBlogDotLogoFunc,
	"@site.cover_image": atBlogDotCoverFunc,
	"@site.description": atBlogDotDescriptionFunc,
	"@site.navigation":  navigationFunc,

	// Post functions
	"post":       postFunc,
	"excerpt":    excerptFunc,
	"title":      titleFunc,
	"content":    contentFunc,
	"post_class": post_classFunc,
	"featured":   featuredFunc,
	"id":         idFunc,
	"post.id":    idFunc,

	// Tag functions
	"tag.name": tagDotNameFunc,
	"tag.slug": tagDotSlugFunc,

	// Author functions
	"bio":                     bioFunc,
	"email":                   emailFunc,
	"website":                 websiteFunc,
	"cover":                   coverFunc,
	"location":                locationFunc,
	"primary_author":          authorFunc,
	"primary_author.name":     authorDotNameFunc,
	"primary_author.bio":      bioFunc,
	"primary_author.email":    emailFunc,
	"primary_author.website":  websiteFunc,
	"primary_author.image":    authorDotImageFunc,
	"primary_author.cover":    coverFunc,
	"primary_author.location": locationFunc,
	// TODO: remove legacy API v2
	"author":          authorFunc,
	"author.name":     authorDotNameFunc,
	"author.bio":      bioFunc,
	"author.email":    emailFunc,
	"author.website":  websiteFunc,
	"author.image":    authorDotImageFunc,
	"author.cover":    coverFunc,
	"author.location": locationFunc,

	// Navigation functions
	"navigation": navigationFunc,
	"label":      labelFunc,
	"current":    currentFunc,
	"slug":       slugFunc,

	// Multiple block functions
	"@first": atFirstFunc,
	"@last":  atLastFunc,
	"@even":  atEvenFunc,
	"@odd":   atOddFunc,
	"name":   nameFunc,
	"url":    urlFunc,

	// Pagination functions
	"pagination": paginationFunc,
	"prev":       prevFunc,
	"next":       nextFunc,
	"page":       pageFunc,
	"pages":      pagesFunc,
	"page_url":   page_urlFunc,
	"pageUrl":    page_urlFunc,

	// Possible if arguments
	"posts":           postsFunc,
	"tags":            tagsFunc,
	"pagination.prev": prevFunc,
	"pagination.next": nextFunc,

	// Possible plural arguments
	"pagination.total":    paginationDotTotalFunc,
	"../pagination.total": paginationDotTotalFunc,
}
