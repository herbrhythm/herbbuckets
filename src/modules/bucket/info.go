package bucket

type Fileinfo struct {
	Size    int64
	Modtime int64
}

func NewFileinfo() *Fileinfo {
	return &Fileinfo{}
}

type DownloadInfo struct {
	URL       string
	ExpiredAt int64
	Permanent bool
}

func NewDownloadInfo() *DownloadInfo {
	return &DownloadInfo{}
}

type CompleteInfo struct {
	ID      string
	Bucket  string
	Object  string
	Size    int64
	Preview *DownloadInfo
}

func NewCompleteInfo() *CompleteInfo {
	return &CompleteInfo{
		Preview: NewDownloadInfo(),
	}
}

type WebuploadInfo struct {
	ID             string
	Bucket         string
	Object         string
	UploadURL      string
	Complete       string
	UploadType     string
	SizeLimit      int64
	ExpiredAt      int64
	PostBody       map[string]string
	FileField      string
	SuccessCodeMin int
	SuccessCodeMax int
}

func NewWebuploadInfo() *WebuploadInfo {
	return &WebuploadInfo{
		PostBody: map[string]string{},
	}
}
