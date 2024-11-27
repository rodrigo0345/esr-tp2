package bootstrapsystem

import "errors"

func (bs *BootstrapSystem) ProcessNeighbors(data *ChannelData) {

  bs.Logger.Debug("Processing neighbors")
	header := data.Msg

	bs.Neighbors = header.GetBootstraperResult().Neighbors

	if bs.Neighbors == nil {
		bs.Neighbors = []string{}
		data.Callback <- &CallbackData{
			Success: false,
			Error:   errors.New("Neighbors are nil"),
		}
    return
	}

	data.Callback <- &CallbackData{
		Success: true,
		Error:   nil,
	}
}
