package task_service

const NUM_每页搜索结果数 = 10

func HasNextSearchPage(searchedPage int, checkRank int) bool {
	return searchedPage < (checkRank / NUM_每页搜索结果数)
}
