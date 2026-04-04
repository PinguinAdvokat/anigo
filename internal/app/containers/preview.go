package containers

import (
	"anigo/internal/extractors"
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/image/webp"
)

type Preview struct {
	*tview.Flex
	app    controller
	Client *http.Client

	Cover       *tview.Image
	Description *tview.TextView
}

func NewPreview(app controller, client *http.Client) *Preview {
	p := &Preview{
		Flex:        tview.NewFlex(),
		app:         app,
		Client:      client,
		Cover:       tview.NewImage(),
		Description: tview.NewTextView(),
	}
	p.Flex.SetBorder(true)
	p.Flex.SetTitle("Описание")
	p.Flex.SetBorderPadding(1, 0, 0, 0)
	p.Flex.Clear()
	p.Flex.SetDirection(tview.FlexColumn).
		AddItem(p.Cover, 0, 2, false).
		AddItem(p.Description, 0, 1, false)

	p.SetBoxResizeFunc(func() {
		p.Flex.ResizeItem(p.Cover, p.getCoverWidth(), 1)
	})
	p.SetTitleColor(tcell.ColorYellow)
	return p
}

func (p *Preview) SetPreview(anime *extractors.Anime) {
	p.SetTitle(anime.Title)
	if anime.Description != "" {
		p.Description.SetText(anime.Description)
	} else {
		p.Description.SetText("нет описания")
	}
	p.Clear().
		AddItem(nil, p.getCoverWidth(), 1, false).
		AddItem(p.Description, 0, 1, false)

	go func() {
		if anime.CoverURL != "" {
			p.SetImageURL(anime.CoverURL)
			p.Clear().
				AddItem(p.Cover, p.getCoverWidth(), 1, false).
				AddItem(p.Description, 0, 1, false)
		} else {
			p.Clear().
				AddItem(tview.NewTextView().SetText("нет фото").SetTextAlign(tview.AlignCenter), 10, 1, false).
				AddItem(p.Description, 0, 1, false)
		}
	}()
	if anime.CoverURL != "" {
		go func() {
			p.SetImageURL(anime.CoverURL)
		}()
	}
}

func (p *Preview) SetSpinner() {
	p.Clear()
	p.AddItem(p.app.GetSpinner(), 0, 1, false)
}

func (p *Preview) getCoverWidth() int {
	_, _, _, height := p.Cover.GetRect()
	return max(int(float32(height)*1.5), 4)
}

func (p *Preview) SetImageURL(url string) {
	img, err := p.loadFromURL(url)
	if err != nil {
		log.Printf("error in displaying image: %v", err)
		return
	}
	p.Cover.SetImage(img).SetColors(tview.TrueColor)
}

func (p *Preview) loadFromURL(url string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		return nil, fmt.Errorf("not an image: %s", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var img image.Image

	// WEBP — через webp.Decode
	if strings.HasSuffix(url, ".webp") || strings.Contains(ct, "webp") {
		img, err = webp.Decode(bytes.NewReader(body))
	} else {
		// PNG/JPEG — через стандартный image.Decode
		img, _, err = image.Decode(bytes.NewReader(body))
	}
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}

	return img, nil
}
