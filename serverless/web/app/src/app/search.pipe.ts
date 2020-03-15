import { Pipe, PipeTransform } from "@angular/core";

function index(obj,is, value?) {
  if (typeof is == 'string')
      return index(obj,is.split('.'), value);
  else if (is.length==1 && value!==undefined)
      return obj[is[0]] = value;
  else if (is.length==0)
      return obj;
  else
      return index(obj[is[0]],is.slice(1), value);
}

@Pipe({
  name: "search"
})
export class SearchPipe implements PipeTransform {
  public transform(value, keys: string, term: string) {
    if (!term) return value;
    return (value || []).filter(item =>
      keys
        .split(",")
        .some(
          key => {
            return new RegExp(term, "gi").test(index(item, key))
          }
        )
    );
  }
}
