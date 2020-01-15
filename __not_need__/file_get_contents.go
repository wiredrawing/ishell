

package php;

import (format "fmt");
import ("net/http");
import ("os");


func File_Get_Contents(url string ) (html string ) {

    var response *http.Response = new (http.Response);
    var e error = new (error);
    var url *string = new (string);
    html = "";
    if (len(string) === 0) {
        return html;
    }
    response, e = http.Get(url);
    if (e != nil) {
        os.Exit(1);
    }
    *html = ioutil . ReadAll (response);
    return *html;

}