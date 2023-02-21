import { BASEPATH } from "@/config";
import { SearchResultApi, SearchResultGrouped } from "@/types/SearchResult";
import axios from "axios";

abstract class SearchService {
  static fetchResults = async (
    query: string
  ): Promise<SearchResultGrouped[]> => {
    const url = `${BASEPATH}/search?q=${encodeURIComponent(query)}`;
    let response = axios.get<SearchResultApi[]>(url).then((res) => {
      const data = res.data;

      let play_names = data
        .map((x) => x.play_name)
        .filter((value, index, array) => array.indexOf(value) === index);

      let responseData: SearchResultGrouped[] = [];

      play_names.forEach((play) => {
        let samePlayItems = data.filter((x) => x.play_name === play);
        responseData.push({
          play_name: play,
          group: samePlayItems,
        });
      });

      return responseData;
    });

    return response;
  };
}

export default SearchService;
