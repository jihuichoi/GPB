package main

import (
	"errors"
)

// ErrNoAvatarURL ErrNoAvatar is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specificed client
	// or returns an error if something goes wrong.
	// ErrNOAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar is...
type AuthAvatar struct{}

// UseAuthAvatar is ...
var UseAuthAvatar AuthAvatar

// GetAvatarURL is ...
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar is ...
type GravatarAvatar struct{}

// UseGravatar is ...
var UseGravatar GravatarAvatar

// GetAvatarURL is ...
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	// if email, ok := c.userData["email"]; ok {
	// 	if emailStr, ok := email.(string); ok {
	// 		m := md5.New()
	// 		io.WriteString(m, strings.ToLower(emailStr))
	// 		return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
	// 	}
	// }
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
