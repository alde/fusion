package db

import (
	"time"

	"github.com/kisielk/sqlstruct"
)

const newsSQL = `SELECT tn.*, tc.*, user_id, user_name FROM fusion_news tn
                LEFT JOIN fusion_users tu ON tn.news_name=tu.user_id
                LEFT JOIN fusion_news_cats tc ON tn.news_cat=tc.news_cat_id
                WHERE ?
									AND (news_start='0'||news_start<=?)
									AND (news_end='0'||news_end>=?)
									AND news_draft='0'
                ORDER BY news_sticky DESC, news_datestamp DESC LIMIT ?,?`

// News represents a news entry
type News struct {
	ID            int    `sql:"news_id" json:"id"`
	Subject       string `sql:"news_subject" json:"subject"`
	Category      int    `sql:"news_cat" json:"category"`
	Summary       string `sql:"news_news" json:"summary"`
	FullText      string `sql:"news_extended" json:"full_text"`
	Breaks        bool   `sql:"news_breaks" json:"breaks"`
	Name          int    `sql:"news_name" json:"name"`
	DateStamp     int    `sql:"news_datestamp" json:"datestamp"`
	Start         int    `sql:"news_start" json:"start"`
	End           int    `sql:"news_end" json:"end"`
	Visibility    int    `sql:"news_visibility" json:"visibility"`
	Reads         int    `sql:"news_reads" json:"reads"`
	Draft         bool   `sql:"news_draft" json:"draft"`
	Sticky        bool   `sql:"news_sticky" json:"sticky"`
	AllowComments bool   `sql:"news_allow_comment" json:"allow_comments"`
	AllowRatings  bool   `sql:"news_allow_ratings" json:"allow_ratings"`
	CategoryID    int    `sql:"news_cat_id" json:"category_id"`
	CategoryName  string `sql:"news_cat_name" json:"category_name"`
	CategoryImage string `sql:"news_cat_image" json:"category_image"`
	UserID        int    `sql:"user_id" json:"user_id"`
	UserName      string `sql:"user_name" json:"user_name"`
}

// News is a function used to query news from the database
func (f *FusionDAO) News(offset, limit int) ([]News, error) {
	sql, err := f.connect()
	if err != nil {
		return nil, err
	}
	now := time.Now().UnixNano()
	// TODO: 1=1 is a bad check, but useful for development
	rows, err := sql.Query(newsSQL, "1=1", now, now, offset, limit)
	if err != nil {
		return nil, err
	}

	var news []News
	for rows.Next() {
		var newsEntry News
		sqlstruct.Scan(&newsEntry, rows)
		news = append(news, newsEntry)
	}
	return news, nil
}
