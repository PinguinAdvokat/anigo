package app

// import (
// 	"fmt"
// 	"log"
// 	"runtime"

// 	"github.com/rivo/tview"
// )

// func (a *App) GetAnimeInfo(index int) {
// 	log.Printf("gorutines: %d\n", runtime.NumGoroutine())
// 	if index < 0 || index > len(a.Manager.FoundAnime)-1 {
// 		log.Printf("index %d out of range in GetAnimeInfo", index)
// 		return
// 	}

// 	if len(a.Manager.FoundAnime[index].AvailableVoiceover) > 0 {
// 		a.Voiceover.SetOptions(a.Manager.FoundAnime[index].AvailableVoiceover, nil)
// 		a.Player.SetOptions(a.Manager.FoundAnime[index].AvailablePlayers, nil)
// 		a.setAnimeSettingsContent([]tview.Primitive{a.Voiceover, a.Player})
// 		a.Draw()
// 		return
// 	}

// 	a.setAnimeSettingsContent([]tview.Primitive{a.Spinner})
// 	go func() {
// 		err := a.Manager.ParseAnime(index)
// 		if err != nil {
// 			log.Printf("Failed get animeinfo: %v\n", err)
// 			a.ErrorView.SetText(fmt.Sprintf("Ошибка при поиске: %v", err))
// 			a.setAnimeSettingsContent([]tview.Primitive{a.ErrorView})
// 			a.Draw()
// 			return
// 		}
// 		a.Voiceover.SetOptions(a.Manager.FoundAnime[index].AvailableVoiceover, nil)
// 		a.Player.SetOptions(a.Manager.FoundAnime[index].AvailablePlayers, nil)
// 		a.setAnimeSettingsContent([]tview.Primitive{a.Voiceover, a.Player})
// 		a.Draw()
// 	}()
// }

// func (a *App) setAnimeSettingsContent(prims []tview.Primitive) {
// 	a.AnimeSettings.Clear()
// 	a.AnimeSettings.SetDirection(tview.FlexRow)
// 	for _, prim := range prims {
// 		a.AnimeSettings.AddItem(prim, 0, 1, true)
// 	}
// }
