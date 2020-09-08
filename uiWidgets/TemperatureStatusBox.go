package uiWidgets

import (
	// "time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type TemperatureStatusBox struct {
	*gtk.Box
	interfaces.ITemperatureDataDisplay

	client						*octoprint.Client
	labelWithImages				map[string]*utils.LabelWithImage
}

func CreateTemperatureStatusBox(
	client						*octoprint.Client,
	includeHotends				bool,
	includeBed					bool,
) *TemperatureStatusBox {
	if !includeHotends && !includeBed {
		utils.Logger.Error("TemperatureStatusBox.CreateTemperatureStatusBox() - both includeToolheads and includeBed are false, but at least one needs to be true")
		return nil
	}

	currentTemperatureData, err := utils.GetCurrentTemperatureData(client)
	if err != nil {
		utils.LogError("TemperatureStatusBox.CreateTemperatureStatusBox()", "GetCurrentTemperatureData(client)", err)
		return nil
	}

	base := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)

	instance := &TemperatureStatusBox{
		Box:						base,
		client:						client,
		labelWithImages:			map[string]*utils.LabelWithImage{},
	}

	instance.SetVAlign(gtk.ALIGN_CENTER)
	instance.SetHAlign(gtk.ALIGN_CENTER)

	var bedTemperatureData *octoprint.TemperatureData = nil
	var hotendIndex int = 0
	var hotendCount int = utils.GetToolheadCount(client)
	for key, temperatureData := range currentTemperatureData {
		if key == "bed" {
			bedTemperatureData = &temperatureData
		} else {
			hotendIndex++

			if includeHotends {
				strImageFileName := utils.GetToolheadFileName(hotendIndex, hotendCount)
				instance.labelWithImages[key] = utils.MustLabelWithImage(strImageFileName, "")
				instance.Add(instance.labelWithImages[key])
			}
		}
	}

	if bedTemperatureData != nil {
		if includeBed {
			instance.labelWithImages["bed"] = utils.MustLabelWithImage("bed.svg", "")
			instance.Add(instance.labelWithImages["bed"])
		}
	}

	if utils.UpdateTemperaturesBackgroundTask == nil {
		utils.CreateUpdateTemperaturesBackgroundTask(instance, client)
	} else {
		utils.RegisterTemperatureStatusBox(instance, client)
	}

	return instance
}

// interfaces.ITemperatureDataDisplay
func (this *TemperatureStatusBox) UpdateTemperatureData(currentTemperatureData map[string]octoprint.TemperatureData) {
	for key, temperatureData := range currentTemperatureData {
		if labelWithImage, ok := this.labelWithImages[key]; ok {
			temperatureDataString := utils.GetTemperatureDataString(temperatureData)
			labelWithImage.Label.SetText(temperatureDataString)
			labelWithImage.ShowAll()
		}
	}
}
