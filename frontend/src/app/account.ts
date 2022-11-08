export interface Account {
  // @JsonProperty('first-line')
  accountId: string;
  typeId: string;
  name: string;
  description: string;
  active: boolean;
}
