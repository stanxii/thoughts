package post

import (
	"errors"

	pb "postservice/genproto"
)

// DbClient manages the communication between services and database
type DbClient struct {
	db *Database
}

// NewDbClient creates a new DbClient instance
func NewDbClient(db *Database) *DbClient {
	return &DbClient{db}
}

// CreatePost creates a new Post row in the database
func (c *DbClient) CreatePost(content string, userID int32) (pb.Post, error) {
	conn := c.db.GetConn()

	var (
		id          int32
		dateCreated string
	)

	err := conn.QueryRow(`INSERT INTO thoughts.posts(content, user_id) VALUES($1, $2)
    RETURNING id, content, user_id, time_format(date_created)`,
		content, userID).Scan(&id, &content, &userID, &dateCreated)
	if err != nil {
		return pb.Post{}, errors.New("Error happened while saving to the database")
	}

	post := pb.Post{
		Id:          id,
		Content:     content,
		UserId:      userID,
		DateCreated: dateCreated}
	return post, nil
}

// GetPost returns a Post row with the passed id from the database
func (c *DbClient) GetPost(id int32) (pb.Post, error) {
	conn := c.db.GetConn()

	var (
		content     string
		userID      int32
		dateCreated string
	)

	err := conn.QueryRow(`SELECT id, content, user_id, time_format(date_created)
    FROM thoughts.posts WHERE id = $1`, id).Scan(&id, &content, &userID, &dateCreated)
	if err != nil {
		return pb.Post{}, errors.New("Error happened while reading from the database")
	}

	post := pb.Post{
		Id:          id,
		Content:     content,
		UserId:      userID,
		DateCreated: dateCreated}
	return post, nil
}

// GetFeed returns posts and retweets of users followed by the user
func (c *DbClient) GetFeed(userID int32, page int32, limit int32) (pb.Posts, error) {
	conn := c.db.GetConn()

	rows, err := conn.Query(`(SELECT id, content, user_id, time_format(date_created)
	FROM thoughts.posts  WHERE user_id = $1 ORDER BY date_created DESC)
	UNION
	(SELECT id, content, user_id, time_format(date_created) FROM thoughts.posts 
	WHERE id IN (SELECT post_id FROM thoughts.retweets
	WHERE user_id = $1 ORDER BY date_created DESC))
	UNION
	(SELECT id, content, user_id, time_format(date_created) FROM thoughts.posts 
	WHERE user_id IN (SELECT user_id FROM thoughts.followings
	WHERE follower_id = $1) ORDER BY date_created DESC)
	OFFSET $2 LIMIT $3`,
		userID, page*limit, limit)
	if err != nil {
		return pb.Posts{}, errors.New("Error happened while reading from the database")
	}
	defer rows.Close()

	var posts []*pb.Post
	for rows.Next() {
		var (
			id          int32
			content     string
			userID      int32
			dateCreated string
		)
		if err := rows.Scan(&id, &content, &userID, &dateCreated); err != nil {
			return pb.Posts{}, errors.New("Error happened while parsing database response")
		}
		post := pb.Post{
			Id:          id,
			Content:     content,
			UserId:      userID,
			DateCreated: dateCreated}
		posts = append(posts, &post)
	}
	return pb.Posts{Posts: posts}, nil
}

// GetPostsCount returns the number of posts and retweets of the user with userID
func (c *DbClient) GetPostsCount(userID int32) (int32, error) {
	conn := c.db.GetConn()

	var count int32

	err := conn.QueryRow(`SELECT COUNT(*) FROM thoughts.posts
    WHERE user_id = $1 OR id IN (SELECT post_id
    FROM thoughts.retweets WHERE user_id = $1)`,
		userID).Scan(&count)
	if err != nil {
		return 0, errors.New("Error happened while reading from the database")
	}
	return count, nil
}

// GetPosts returns the posts and retweets of the user with userID
func (c *DbClient) GetPosts(userID int32, page int32, limit int32) (pb.Posts, error) {
	conn := c.db.GetConn()

	rows, err := conn.Query(`(SELECT id, content, user_id, time_format(date_created)
	FROM thoughts.posts  WHERE user_id = $1 ORDER BY date_created DESC)
	UNION
	(SELECT id, content, user_id, time_format(date_created) FROM thoughts.posts 
	WHERE id IN (SELECT post_id FROM thoughts.retweets
	WHERE user_id = $1 ORDER BY date_created DESC))
	OFFSET $2 LIMIT $3`,
		userID, page*limit, limit)

	if err != nil {
		return pb.Posts{}, errors.New("Error happened while reading from the database")
	}
	defer rows.Close()

	var posts []*pb.Post
	for rows.Next() {
		var (
			id          int32
			content     string
			userID      int32
			dateCreated string
		)
		if err := rows.Scan(&id, &content, &userID, &dateCreated); err != nil {
			return pb.Posts{}, errors.New("Error happened while parsing database response")
		}
		post := pb.Post{
			Id:          id,
			Content:     content,
			UserId:      userID,
			DateCreated: dateCreated}
		posts = append(posts, &post)
	}
	return pb.Posts{Posts: posts}, nil
}

// GetLikedPostsCount returns the number of liked posts of the user with userID
func (c *DbClient) GetLikedPostsCount(userID int32) (int32, error) {
	conn := c.db.GetConn()

	var count int32

	err := conn.QueryRow(`SELECT COUNT(*) FROM thoughts.posts
    WHERE id IN (SELECT post_id FROM thoughts.likes
    WHERE user_id = $1)`, userID).Scan(&count)
	if err != nil {
		return 0, errors.New("Error happened while reading from the database")
	}
	return count, nil
}

// GetLikedPosts returns posts liked by the user
func (c *DbClient) GetLikedPosts(userID int32, page int32, limit int32) (pb.Posts, error) {
	conn := c.db.GetConn()

	rows, err := conn.Query(`SELECT id, content, user_id, time_format(date_created)
    FROM thoughts.posts
    WHERE id IN (SELECT post_id FROM thoughts.likes WHERE user_id = $1
    ORDER BY date_created DESC) OFFSET $2 LIMIT $3`,
		userID, page*limit, limit)
	if err != nil {
		return pb.Posts{}, errors.New("Error happened while reading from the database")
	}
	defer rows.Close()

	var posts []*pb.Post
	for rows.Next() {
		var (
			id          int32
			content     string
			userID      int32
			dateCreated string
		)
		if err := rows.Scan(&id, &content, &userID, &dateCreated); err != nil {
			return pb.Posts{}, errors.New("Error happened while parsing database response")
		}
		post := pb.Post{
			Id:          id,
			Content:     content,
			UserId:      userID,
			DateCreated: dateCreated}
		posts = append(posts, &post)
	}
	return pb.Posts{Posts: posts}, nil
}

// DeletePost deletes a post with the passed postID
func (c *DbClient) DeletePost(postID int32) error {
	conn := c.db.GetConn()

	_, err := conn.Exec(`DELETE FROM thoughts.posts WHERE id = $1`, postID)
	if err != nil {
		return errors.New("Error happened while writing to the database")
	}

	return nil
}

// LikePost creates a like relationship between user and post
func (c *DbClient) LikePost(postID int32, userID int32) error {
	conn := c.db.GetConn()

	_, err := conn.Exec(`INSERT INTO thoughts.likes (post_id, user_id)
    VALUES ($1, $2)`, postID, userID)
	if err != nil {
		return errors.New("Error happened while writing to the database")
	}

	return nil
}

// UnlikePost deletes a like relationship between user and post
func (c *DbClient) UnlikePost(postID int32, userID int32) error {
	conn := c.db.GetConn()

	_, err := conn.Exec(`DELETE FROM thoughts.likes
    WHERE post_id = $1 AND user_id = $2`, postID, userID)
	if err != nil {
		return errors.New("Error happened while writing to the database")
	}

	return nil
}

// RetweetPost creates a retweet relationship between user and post
func (c *DbClient) RetweetPost(postID int32, userID int32) error {
	conn := c.db.GetConn()

	_, err := conn.Exec(`INSERT INTO thoughts.retweets (post_id, user_id)
    VALUES ($1, $2)`, postID, userID)
	if err != nil {
		return errors.New("Error happened while writing to the database")
	}

	return nil
}

// RemoveRetweet deletes a retweet relationship between user and post
func (c *DbClient) RemoveRetweet(postID int32, userID int32) error {
	conn := c.db.GetConn()

	_, err := conn.Exec(`DELETE FROM thoughts.retweets
    WHERE post_id = $1 AND user_id = $2`, postID, userID)
	if err != nil {
		return errors.New("Error happened while writing to the database")
	}

	return nil
}
