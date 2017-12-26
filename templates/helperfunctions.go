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

	// @blog functions
	"@blog.title":       atBlogDotTitleFunc,
	"@blog.url":         atBlogDotUrlFunc,
	"@blog.logo":        atBlogDotLogoFunc,
	"@blog.cover":       atBlogDotCoverFunc,
	"@blog.description": atBlogDotDescriptionFunc,
	"@blog.navigation":  navigationFunc,

	// Post functions
	"post":       postFunc,
	"excerpt":    excerptFunc,
	"title":      titleFunc,
	"content":    contentFunc,
	"post_class": post_classFunc,
	"featured":   featuredFunc,
	"published":  publishedFunc,
	"id":         idFunc,
	"post.id":    idFunc,
	"prev_post":  prevPostFunc,
	"next_post":  nextPostFunc,

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
	"twitter":         twitterFunc,
	"facebook":        facebookFunc,
	"author.name":     authorDotNameFunc,
	"author.bio":      bioFunc,
	"author.email":    emailFunc,
	"author.website":  websiteFunc,
	"author.image":    authorDotImageFunc,
	"author.cover":    coverFunc,
	"author.location": locationFunc,
	"author.twitter":  twitterFunc,
	"author.facebook": facebookFunc,

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
