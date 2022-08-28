package archive

import (
	"fmt"
	"sync"
	"testing"
)

// func TestArchive(t *testing.T) {
// 	if err := run(); err != nil {
// 		t.Fatal(err)
// 	}
// 	// t.Fatal("success")
// }

// func run() (err error) {
// 	archive, err := NewArchive("test.db")
// 	if err != nil {
// 		return
// 	}

// 	const firstUrl = "https://golang.org"

// 	var result bool
// 	result, err = archive.IsLoaded(firstUrl)
// 	if err != nil {
// 		return
// 	}

// 	if result != false {
// 		return fmt.Errorf("expected false for unknown item, got %v", result)
// 	}

// 	err = archive.SetLoaded(firstUrl, true)
// 	if err != nil {
// 		return
// 	}

// 	err = archive.SetLoaded("https://another.org", true)
// 	if err != nil {
// 		return
// 	}

// 	result, err = archive.IsLoaded(firstUrl)
// 	if err != nil {
// 		return
// 	}

// 	if result != true {
// 		return fmt.Errorf("expected true for known item, got %v", result)
// 	}

// 	return
// }

// func TestAppend(t *testing.T) {
// 	if err := Append(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func Append() (err error) {
// 	const appendFilename = "append-test.db"
// 	archive, err := NewRepo(appendFilename)
// 	if err != nil {
// 		return
// 	}

// 	const dirpath = "./testdata"
// 	dir, err := os.ReadDir(dirpath)
// 	if err != nil {
// 		return
// 	}

// 	for _, entry := range dir {
// 		fullpath := filepath.Join(dirpath, entry.Name())
// 		fmt.Println(fullpath)
// 		var bytes []byte
// 		bytes, err = os.ReadFile(fullpath)
// 		if err != nil {
// 			return
// 		}

// 		var item Item
// 		err = json.Unmarshal(bytes, &item)
// 		if err != nil {
// 			return
// 		}

// 		var archiveItem ArchiveItem
// 		archiveItem.HistoryKey = item.Basename
// 		archiveItem.Description = &ArchiveItemDescription{
// 			Date:      item.Date,
// 			Source:    item.Provider,
// 			ShowName:  item.ShowName,
// 			ShowTitle: item.ShowTitle,
// 			EpTitle:   item.EpTitle,
// 			Artists:   item.Artists,
// 		}

// 		archiveItem.Meta = &ArchiveItemMeta{
// 			Filename: item.Filename,
// 			Duration: item.Duration,
// 			Size:     item.Size,
// 		}

// 		if item.MessageId != nil {
// 			archiveItem.Chan = &ArchiveItemChan{
// 				MessageId: *item.MessageId,
// 			}
// 		}

// 		fmt.Println(archiveItem)

// 		err = archive.Add(archiveItem)
// 		if err != nil {
// 			return
// 		}
// 	}

// 	return nil
// }

func TestConcurrentReadWrite(t *testing.T) {
	const filename = "concurrent-test.db"

	archive, err := CreateDB(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	archiveRepo, err := NewRepo()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := archiveRepo.Migrate(archive); err != nil {
		fmt.Println(err)
		return
	}

	const iter = 10
	wg := sync.WaitGroup{}
	wg.Add(iter)
	wg.Add(iter)

	for i := 0; i < iter; i++ {
		key := fmt.Sprint(i)
		go func() {
			err := archiveRepo.Create(archive, key, true, ArchiveItem{HistoryKey: key})
			if err != nil {
				fmt.Printf("key(%s) create error: %s\n", key, err.Error())
			}
			wg.Done()
		}()

		go func() {
			ok, err := archiveRepo.IsLoaded(archive, key)
			if err != nil {
				fmt.Printf("key(%s) isLoaded error: %s\n", key, err.Error())
				return
			}

			fmt.Printf("key(%s): presented=%v\n", key, ok)
			wg.Done()
		}()
	}

	wg.Wait()

	for i := 0; i < iter; i++ {
		key := fmt.Sprint(i)
		ok, err := archiveRepo.IsLoaded(archive, key)
		if err != nil {
			fmt.Printf("key(%s) isLoaded error: %s\n", key, err.Error())
			return
		}

		fmt.Printf("key(%s): presented=%v\n", key, ok)
	}

	fmt.Println("done")
}
