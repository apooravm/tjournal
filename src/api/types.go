package api

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateJournalLogReq struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Log      string   `json:"log"`
	Tags     []string `json:"tags"`
	Title    string   `json:"title"`
}

type ReadJournalLogRes struct {
	Created_at string   `json:"created_at"`
	Log        string   `json:"log_message"`
	Title      string   `json:"title"`
	Tags       []string `json:"tags"`
	Log_Id     int      `json:"log_id"`
}

type UpdateLogReq struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Log      string   `json:"log"`
	Tags     []string `json:"tags"`
	Title    string   `json:"title"`
	Log_Id   int      `json:"log_id"`
}

type DeleteJournalLogReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Log_Id   int    `json:"log_id"`
}

type JournalMessage struct {
	Message string
	Code    int
	Simple  string
}
