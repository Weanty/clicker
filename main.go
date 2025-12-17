package main

import (
	//"strconv"
	"fmt"
	"math"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	
	"encoding/json"
  "os"
  "path/filepath"
)


type GameState struct {
	Count 	float64 	`json:"count"`
	Gain 		float64 	`json:"gain"`
	ShopPrices []float64 `json:"shop_prices"`
	UpgradesBought []int `json:"upgrades_bought"`
}

func main() {
	ternary := func(condition bool, trueVal, falseVal float64) float64 {
		if condition {
			return trueVal
		}
		return falseVal
	}

	var count float64 = 0
	var gain float64 = 1
	//var currentIndex int = 0
	const separator float64 = 1000
	var pC1 float64 = 1.3
	var pC3 float64 = 0.4

	conTypes := []string{
		"",
		"thousand",
		"million",
		"billion",
		"trillion",
		"quadrillion",
		"quintillion",
		"sextillion",
		"septillion",
		"octillion",
		"nonillion",
		"decillion",
		"undecillion",
		"duodecillion",
		"tredecillion",
		"quattuordecillion",
		"quindecillion",
		"sexdecillion",
		"septendecillion",
		"octadecillion",
		"novemdecillion",
		"vigintillion",
		"unvigintillion",
		"duovigintillion",
		"trevigintillion",
	}

	shopPrices := []float64{
		100,
		1000,
		math.Pow(separator, float64(2)),
		math.Pow(separator, float64(3)),
	}

	upgradesBought := []int{
		0,
		0,
		0,
		0,
	}

	upgradeMultiplier := []float64{
		1.20,
		1.30,
		1.45,
		1.65,
	}

	application := app.NewWithID("com.weanty.clickergame")
	clickerWindow := application.NewWindow("Clicking Inc.")
	clickerWindow.SetMaster()

	if loadedState, err := loadGameFromFile(application); err == nil {
		count = loadedState.Count
		shopPrices = loadedState.ShopPrices
		upgradesBought = loadedState.UpgradesBought
		gain = loadedState.Gain
	}
	
	var countButton *widget.Button
	
	var upgradeContainers []*fyne.Container
	var uNameList []*widget.Label
	var uCostList []*widget.Label
	var uBuyList []*widget.Button

	var tabs *container.AppTabs
	//var shopTabs *container.AppTabs

	uNames := []string{
		"Stronger Fingers",
		"Iron Fists",
		"Better Mouse",
		"Krisz's DC Server",
	}

	convType := func(cc float64) (newC float64, cType string) {
		for i := 0; i < len(conTypes); i++ {
			if cc < math.Pow(separator, float64(i + 1)) {
				cType = conTypes[i]
				newC = cc / math.Pow(separator, float64(i))
				return newC, cType
			}
		}

		lastI := len(conTypes) - 1
		return cc / math.Pow(separator, float64(lastI)), conTypes[lastI]
	}

	buyUpgrade := func(upgradeId int) {
		if count >= shopPrices[upgradeId] {
			count -= shopPrices[upgradeId]
			upgradesBought[upgradeId]++
			shopPrices[upgradeId] = ternary(upgradesBought[upgradeId] > 4, shopPrices[upgradeId] * (pC1 * float64(upgradesBought[upgradeId]) * pC3), shopPrices[upgradeId] * 1.25)
			gain = gain * upgradeMultiplier[upgradeId]
			
			dispCount, countType := convType(count)
			dispGain, gainType := convType(gain)

			countButton.SetText(fmt.Sprintf("%.2f %s\n%.2f %s gpc", dispCount, countType, dispGain, gainType))
			
			dispCost, costType := convType(shopPrices[upgradeId])

			uCostList[upgradeId].SetText(fmt.Sprintf("%.2f %s", dispCost, costType))
			uBuyList[upgradeId].SetText(fmt.Sprintf("Buy | Bought: %d", upgradesBought[upgradeId]))
		}
	}

	countButton = widget.NewButton("Start Clicking!", func() {
		count += gain
		
		dispCount, countType := convType(count)
		dispGain, gainType := convType(gain)

		countButton.SetText(fmt.Sprintf("%.2f %s\n%.2f %s gpc", dispCount, countType, dispGain, gainType))
	})

	//setting up everything
	for i := 0; i < len(uNames); i++ {
		index := i
		
		dispCost, costType := convType(shopPrices[index])

		nameLabel := widget.NewLabel(uNames[index])
		costLabel := widget.NewLabel(fmt.Sprintf("%.2f %s", dispCost, costType))
		buyButton := widget.NewButton(fmt.Sprintf("Buy | Bought: %d", upgradesBought[index]), func(){
			buyUpgrade(index)
		})

		upgradeContainer := container.NewBorder(
			nil,
			buyButton,
			nameLabel,
			costLabel,
			nil,
		)

		uNameList = append(uNameList, nameLabel)
		uCostList = append(uCostList, costLabel)
		uBuyList = append(uBuyList, buyButton)
		upgradeContainers = append(upgradeContainers, upgradeContainer)
	}

	shopContentArea := container.NewMax()
	
	shopList := widget.NewList(
		func() int { return len(upgradeContainers) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("#%d", i + 1))
		},
	)
	
	shopList.OnSelected = func(id widget.ListItemID) {
		shopContentArea.Objects = []fyne.CanvasObject{upgradeContainers[id]}
		shopContentArea.Refresh()
	}

	shopList.Select(0)
	
	//Setting content
	clickerContent := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		countButton,
	)
	
	saveButton := widget.NewButton("Save Game", func() {
    state := GameState{
        Count:          count,
        Gain:           gain,
        ShopPrices:     shopPrices,
        UpgradesBought: upgradesBought,
    }
    err := saveGameToFile(state, application)
    if err != nil {
        fmt.Println("Save failed:", err)
    } else {
        fmt.Println("Game saved!")
    }
	})

	shopContent := container.NewBorder(nil, nil, shopList, nil, shopContentArea)
	
	settingsTabs := container.NewAppTabs(
    container.NewTabItem("Save", saveButton),
	)

	settingsTabs.SetTabLocation(container.TabLocationLeading)
	
	tabs = container.NewAppTabs(
		container.NewTabItem("Home", clickerContent),
		container.NewTabItem("Shop", shopContent),
		container.NewTabItem("Settings", settingsTabs),
	)
	
	//Initializing windows
	clickerWindow.SetContent(tabs)
	clickerWindow.Resize(fyne.NewSize(300, 100))
	clickerWindow.Show()

	application.Run()
}

func loadGameFromFile(a fyne.App) (GameState, error) {
	savePath := filepath.Join(a.Storage().RootURI().Path(), "savefile.json")

	data, err := os.ReadFile(savePath)
	if err != nil {
		return GameState{}, err
	}

	var state GameState
	err = json.Unmarshal(data, &state)
	return state, err
}

func saveGameToFile(state GameState, a fyne.App) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	savePath := filepath.Join(a.Storage().RootURI().Path(), "savefile.json")
	return os.WriteFile(savePath, data, 0644)
}
