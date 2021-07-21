package search

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const Version = "2.0.1"

// SearchInUrls поиск текста
// Возвращает слайс из урлов, где есть совпадение
func SearchInUrls(query string, urls []string) ([]string, error) {
	defer func() {
		if v := recover(); v != nil {
			fmt.Println("SearchInUrls recovered", v)
		}
	}()

	result := []string{}
	var errorSearch error
	var wg sync.WaitGroup

	for _, u := range urls {
		_, err := url.Parse(u)
		if err != nil {
			return nil, WrapSearchError(fmt.Errorf("Error parse url %s :%v", u, err))
		}

		ch := make(chan string)
		chError := make(chan error)
		wg.Add(1)

		go func() {
			defer wg.Done()
			// defer func() {
			// 	if v := recover(); v != nil {
			// 		// log.Fatalln(WrapSearchError(fmt.Errorf("Recover %s :%v", u, v)))
			// 		fmt.Println("recovered", v)
			// 	}
			// }()

			resp, err := http.Get(u)
			if err != nil {
				ch <- ""
				chError <- WrapSearchError(fmt.Errorf("Error get url %s :%v", u, err))
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ch <- ""
				chError <- WrapSearchError(fmt.Errorf("Error get url %s :%v", u, err))
				return
			}
			ch <- string(body)
			chError <- nil
		}()

		body := <-ch
		errorSearch = <-chError

		if strings.Contains(body, query) {
			result = append(result, u)
		}
	}

	wg.Wait()
	return result, errorSearch
}
