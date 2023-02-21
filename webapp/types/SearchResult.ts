export interface SearchResultApi {
  line: string;
  character?: string;
  act_number?: string;
  scene_number?: string;
  play_name: string;
}

export interface SearchResultGrouped {
  play_name: string;
  group: SearchResultApi[];
}
