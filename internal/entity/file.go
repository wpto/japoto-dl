package entity

import "fmt"

func Load(entity *Entity) (err error) {
	body, err := entity.Loader.Raw(entity.URL, entity.Gopts)
	if err != nil {
		err = fmt.Errorf("LoadChunk: %w", err)
		return
	}

	entity.Body = body
	return
}

func IsBodyEmpty(entity *Entity) (err error) {
	if len(entity.Body) == 0 {
		err = fmt.Errorf("Body is empty: %v", entity.Filename)
	}
	return
}

func IsBodyNull(entity *Entity) (err error) {
	if string(entity.Body) == "null" {
		err = fmt.Errorf(`Body is "null": %v`, entity.Filename)
	}
	return
}

func Save(entity *Entity) (err error) {
	err = entity.Workdir.SaveRaw(entity.Filename, entity.Body)
	if err != nil {
		err = fmt.Errorf("Save file: %w", err)
	}
	return
}

func SaveImage(entity *Entity) (err error) {
	err = entity.Workdir.SaveNamedRaw("image", entity.Body)
	if err != nil {
		err = fmt.Errorf("Save image: %w", err)
	}
	return
}

func DownloadFile(entity *Entity) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Download file: %w", err)
		}
	}()

	url, err := entity.ModelFile.Url(entity.TSAudioURL)
	if err != nil {
		return
	}

	entity.URL = url

	// f := &entity.Entity{
	// 	Type:     entity.FileEntity,
	// 	Gopts:    gopts,
	// 	Loader:   loader,
	// 	URL:      url,
	// 	Filename: file.Name(),
	// }

	err = Load(entity)
	if err != nil {
		return
	}

	entity.ModelFile.SetBody(entity.Body)

	err = IsBodyEmpty(entity)
	if err != nil {
		return
	}

	err = IsBodyNull(entity)
	if err != nil {
		return
	}

	err = Save(entity)
	if err != nil {
		return
	}

	return
}

func (entity *Entity) DownloadImage() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Download image: %w", err)
		}
	}()

	err = Load(entity)
	if err != nil {
		return
	}

	err = SaveImage(entity)
	if err != nil {
		return
	}

	return
}
