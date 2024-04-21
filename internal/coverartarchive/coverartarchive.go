package coverartarchive

import (
	"fmt"
)

const API_BASE = "https://coverartarchive.org"

func GetAlbumArt(releaseMbid string) string {
	return fmt.Sprintf("%s/release/%s/front-%d.jpg", API_BASE, releaseMbid, 250)
}

// func getAlbumArt(releaseMbid string, size int) string {
// 	// albumArtUrl := fmt.Sprintf("https://placehold.jp/%dx%d.png", size, size)
//
// 	path := fmt.Sprintf("%s/release/%s/front-%d.jpg", API_BASE, releaseMbid, size)
//
// 	return path
//
// 	// response, err := http.Get(path)
// 	// if err != nil {
// 	// 	log.Panic(err)
// 	// 	return albumArtUrl
// 	// }
// 	//
// 	// url := response.Request.URL
// 	// if response.StatusCode == 200 && url != nil {
// 	// 	albumArtUrl = url.String()
// 	// }
// 	//
// 	// return albumArtUrl
// }
