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
	"meta_title":       metaTitleFunc,
	"meta_description": metaDescriptionFunc,
	"ghost_head":       ghostHeadFunc,
	"ghost_foot":       ghostFootFunc,
	"body_class":       bodyClassFunc,
	"plural":           pluralFunc,
	"date":             dateFunc,
	"image":            imageFunc,
	"contentFor":       contentForFunc,
	"block":            blockFunc,

	// @blog functions
	"@blog.title":       atBlogDotTitleFunc,
	"@blog.url":         atBlogDotURLFunc,
	"@blog.logo":        atBlogDotLogoFunc,
	"@blog.cover":       atBlogDotCoverFunc,
	"@blog.description": atBlogDotDescriptionFunc,
	"@blog.navigation":  navigationFunc,

	// Post functions
	"post":       postFunc,
	"excerpt":    excerptFunc,
	"title":      titleFunc,
	"content":    contentFunc,
	"post_class": postClassFunc,
	"featured":   featuredFunc,
	"id":         idFunc,
	"post.id":    idFunc,

	// Tag functions
	"tag.name": tagDotNameFunc,
	"tag.slug": tagDotSlugFunc,

	// Author functions
	"author":          authorFunc,
	"bio":             bioFunc,
	"email":           emailFunc,
	"website":         websiteFunc,
	"cover":           coverFunc,
	"location":        locationFunc,
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
	"page_url":   pageURLFunc,
	"pageUrl":    pageURLFunc,

	// Possible if arguments
	"posts":           postsFunc,
	"tags":            tagsFunc,
	"pagination.prev": prevFunc,
	"pagination.next": nextFunc,

	// Possible plural arguments
	"pagination.total":    paginationDotTotalFunc,
	"../pagination.total": paginationDotTotalFunc,
}
