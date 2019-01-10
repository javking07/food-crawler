package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/javking07/food-crawler/model"
)

const apiLookup = "https://api.fda.gov/food/event.json?search=products.industry_code:23&limit=1"

// Retrieve looks up food data and adds to active food-crawler storage, before returning for downstream response to user
func (a *App) RetrieveCount(name string) (*model.Data, error) {
	data := &model.Data{}

	if a.AppCache != nil {
		// TODO otherwise check Cache
		cacheData, err := a.AppCache.Get([]byte(name))
		if err != nil {
			return nil, err
		}
		data.Data = append(data.Data, cacheData)
		return data, nil

		if a.AppDatabase != nil {
		}

	}

	// TODO Default to looking up data then adding to cache/database
	err := a.GetData(data)

	if err != nil {
		return nil, err
	}
	return &model.Data{}, nil
}

func (a *App) GetData(vessel interface{}) error {
	resp, err := a.AppClient.Get(apiLookup)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &vessel)

	if err != nil {
		return err
	}

	fmt.Printf("%+v", vessel)

	return nil
}
