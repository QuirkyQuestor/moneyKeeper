export interface AccountType {
  // @JsonProperty('first-line')
  typeId: string;
  name: string;
  description: string;
}

// function JsonProperty(name: string) {
//   return function DoJsonProperty(
//     target: any,
//     propertyKey: string,
//     descriptor: PropertyDescriptor
//   ) {
//     descriptor.get = function () {
//       return this.data[name];
//     };
//     descriptor.set = function (value) {
//       this.data[name] = value;
//     };
//   };
// }
