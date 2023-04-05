// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-frame/db/mysql/dao/internal"
)

// internalCmsUploadQueueDao is internal type for wrapping internal DAO implements.
type internalCmsUploadQueueDao = *internal.CmsUploadQueueDao

// cmsUploadQueueDao is the data access object for table cms_upload_queue.
// You can define custom methods on it to extend its functionality as you wish.
type cmsUploadQueueDao struct {
	internalCmsUploadQueueDao
}

var (
	// CmsUploadQueue is globally public accessible object for table cms_upload_queue operations.
	CmsUploadQueue = cmsUploadQueueDao{
		internal.NewCmsUploadQueueDao(),
	}
)

// Fill with you ideas below.
