package mybigcamel

import (
    "fmt"
    "strings"
    "testing"
)

func Test_cache(t *testing.T) {
    SS := "OauthIDAPI"
    
    tmp0 := UnMarshal(SS)
    fmt.Println(tmp0)
    tmp1 := Marshal(tmp0)
    fmt.Println(tmp1)
    
    if SS != tmp1 {
        fmt.Println("false.")
    }
    
    fmt.Println(CapLowercase("IDAPIID"))
    fmt.Println(CapSmallcase("IDAPIID"))
}

func CapLowercase(name string) string {
    list := strings.Split(UnMarshal(name), "_")
    if len(list) == 0 {
        return ""
    }
    
    return list[0] + name[len(list[0]):]
}

func CapSmallcase(name string) string {
    list := strings.Split(UnSmallMarshal(name), "_")
    if len(list) == 0 {
        return ""
    }
    
    return list[0] + name[len(list[0]):]
}

func TestUnMarshal(t *testing.T) {
    type args struct {
        name string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {
            args: args{name: "my_name"},
            want: "myName",
        },
        {
            args: args{name: "my_namE"},
            want: "myNamE",
        },
        {
            args: args{name: "myname"},
            want: "myname",
        },
        {
            args: args{name: "MyName"},
            want: "myName",
        },
        {
            args: args{name: "MyNamePrince"},
            want: "myNamePrince",
        },
        {
            args: args{name: "MyName,prince"},
            want: "myNamePrince",
        },
        {
            args: args{name: "MyName,Prince"},
            want: "myNamePrince",
        },
        {
            args: args{name: "MyName, , - prince"},
            want: "myNamePrince",
        },
        {
            args: args{name: "MyName prince"},
            want: "myNamePrince",
        },
        {
            args: args{name: "_MyName prince"},
            want: "MyNamePrince",
        },
        {
            args: args{name: "MyName prince_"},
            want: "myNamePrince",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := UnSmallMarshal(tt.args.name)
            fmt.Println(">>>>:", got)
            if got != tt.want {
                t.Errorf("UnMarshal() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMarshal(t *testing.T) {
    type args struct {
        name string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {
            args: args{name: "my_name"},
            want: "MyName",
        },
        {
            args: args{name: "my_namE"},
            want: "MyNamE",
        },
        {
            args: args{name: "myname"},
            want: "Myname",
        },
        {
            args: args{name: "MyName"},
            want: "MyName",
        },
        {
            args: args{name: "MyNamePrince"},
            want: "MyNamePrince",
        },
        {
            args: args{name: "MyName,prince"},
            want: "MyNamePrince",
        },
        {
            args: args{name: "MyName,Prince"},
            want: "MyNamePrince",
        },
        {
            args: args{name: "MyName, , - prince"},
            want: "MyNamePrince",
        },
        {
            args: args{name: "MyName prince"},
            want: "myNamePrince",
        },
        {
            args: args{name: "_MyName prince"},
            want: "MyNamePrince",
        },
        {
            args: args{name: "MyName prince_"},
            want: "MyNamePrince",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Marshal(tt.args.name); got != tt.want {
                t.Errorf("Marshal() = %v, want %v", got, tt.want)
            }
        })
    }
}
