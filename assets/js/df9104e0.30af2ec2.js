"use strict";(self.webpackChunkdocumentation=self.webpackChunkdocumentation||[]).push([[448],{8453:(e,n,t)=>{t.d(n,{R:()=>i,x:()=>l});var r=t(6540);const s={},a=r.createContext(s);function i(e){const n=r.useContext(a);return r.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function l(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:i(e.components),r.createElement(a.Provider,{value:n},e.children)}},9346:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>c,contentTitle:()=>l,default:()=>d,frontMatter:()=>i,metadata:()=>r,toc:()=>o});const r=JSON.parse('{"id":"type-handling/special-types","title":"Special Types","description":"Censor provides support for special Go types that require specific handling. This guide details how these types are processed.","source":"@site/docs/type-handling/special-types.md","sourceDirName":"type-handling","slug":"/type-handling/special-types","permalink":"/censor/type-handling/special-types","draft":false,"unlisted":false,"editUrl":"https://github.com/vpakhuchyi/censor/tree/main/documentation/docs/type-handling/special-types.md","tags":[],"version":"current","lastUpdatedBy":"Viktor","lastUpdatedAt":1743858422000,"sidebarPosition":3,"frontMatter":{"sidebar_position":3},"sidebar":"docs","previous":{"title":"Complex Types","permalink":"/censor/type-handling/complex-types"},"next":{"title":"Format-Specific","permalink":"/censor/type-handling/format-specific"}}');var s=t(4848),a=t(8453);const i={sidebar_position:3},l="Special Types",c={},o=[{value:"Custom Types",id:"custom-types",level:2},{value:"Type Aliases",id:"type-aliases",level:3},{value:"Custom Structs",id:"custom-structs",level:3},{value:"Interfaces",id:"interfaces",level:2},{value:"Any Type",id:"any-type",level:2},{value:"Custom Type Handlers",id:"custom-type-handlers",level:2},{value:"Circular References",id:"circular-references",level:2},{value:"Reflection Edge Cases",id:"reflection-edge-cases",level:2},{value:"Zero Values",id:"zero-values",level:3},{value:"Unexported Fields",id:"unexported-fields",level:3},{value:"Complete Example",id:"complete-example",level:2},{value:"Next Steps",id:"next-steps",level:2}];function p(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",...(0,a.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.header,{children:(0,s.jsx)(n.h1,{id:"special-types",children:"Special Types"})}),"\n",(0,s.jsx)(n.p,{children:"Censor provides support for special Go types that require specific handling. This guide details how these types are processed."}),"\n",(0,s.jsx)(n.h2,{id:"custom-types",children:"Custom Types"}),"\n",(0,s.jsx)(n.h3,{id:"type-aliases",children:"Type Aliases"}),"\n",(0,s.jsx)(n.p,{children:"Censor handles type aliases by processing them according to their underlying type:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'// String alias\ntype Email string\n\n// Struct with Email type\ntype User struct {\n    ID    string `censor:"display"`\n    Email Email\n}\n\nuser := User{\n    ID:    "123",\n    Email: "user@example.com",\n}\n\n// TEXT output: {ID:123 Email:[CENSORED]}\n// JSON output: {"ID":"123","Email":"[CENSORED]"}\n'})}),"\n",(0,s.jsx)(n.h3,{id:"custom-structs",children:"Custom Structs"}),"\n",(0,s.jsx)(n.p,{children:"Custom struct types are processed like regular structs:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'// Custom struct type\ntype Credentials struct {\n    Username string\n    Password string\n}\n\n// Type alias for Credentials\ntype UserCredentials Credentials\n\n// Usage\ncreds := UserCredentials{\n    Username: "johndoe",\n    Password: "secret123",\n}\n\n// TEXT output: {Username:[CENSORED] Password:[CENSORED]}\n// JSON output: {"Username":"[CENSORED]","Password":"[CENSORED]"}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"interfaces",children:"Interfaces"}),"\n",(0,s.jsx)(n.p,{children:"Censor processes interfaces by examining their concrete type at runtime and applying the appropriate masking:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'// Interface handling\ntype Data interface{}\n\n// Use with different types\nvar stringData Data = "sensitive-data"\nvar intData Data = 42\nvar boolData Data = true\n\n// String will be masked\n// TEXT output: [CENSORED]\n// JSON output: "[CENSORED]"\n\n// Int will be displayed\n// TEXT output: 42\n// JSON output: 42\n\n// Bool will be displayed\n// TEXT output: true\n// JSON output: true\n\n// Struct will have its fields masked according to the rules\ntype User struct {\n    ID    string `censor:"display"`\n    Email string\n}\n\nvar userData Data = User{\n    ID:    "123",\n    Email: "user@example.com",\n}\n\n// TEXT output: {ID:123 Email:[CENSORED]}\n// JSON output: {"ID":"123","Email":"[CENSORED]"}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"any-type",children:"Any Type"}),"\n",(0,s.jsxs)(n.p,{children:["The ",(0,s.jsx)(n.code,{children:"any"})," type (alias for ",(0,s.jsx)(n.code,{children:"interface{}"}),") is handled the same way as interfaces:"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'// Map with any values\ndata := map[string]any{\n    "id":     "123",\n    "email":  "user@example.com",\n    "active": true,\n    "score":  95,\n    "metadata": map[string]string{\n        "lastLogin": "2023-01-01",\n    },\n}\n\n// TEXT output: map[id:123 email:[CENSORED] active:true score:95 metadata:map[lastLogin:[CENSORED]]]\n// JSON output: {"id":"123","email":"[CENSORED]","active":true,"score":95,"metadata":{"lastLogin":"[CENSORED]"}}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"custom-type-handlers",children:"Custom Type Handlers"}),"\n",(0,s.jsx)(n.p,{children:"You can register custom handlers for specific types:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'package main\n\nimport (\n    "fmt"\n    "time"\n    "github.com/vpakhuchyi/censor"\n)\n\n// Custom type\ntype CreditCard struct {\n    Number   string\n    ExpMonth int\n    ExpYear  int\n    CVV      string\n}\n\nfunc main() {\n    // Register custom type handler\n    c := censor.New(\n        censor.WithTypeHandler(func(cc CreditCard) interface{} {\n            // Mask all except last 4 digits\n            if len(cc.Number) > 4 {\n                return CreditCard{\n                    Number:   "****-****-****-" + cc.Number[len(cc.Number)-4:],\n                    ExpMonth: cc.ExpMonth,\n                    ExpYear:  cc.ExpYear,\n                    CVV:      "***",\n                }\n            }\n            return CreditCard{\n                Number:   "[CENSORED]",\n                ExpMonth: cc.ExpMonth,\n                ExpYear:  cc.ExpYear,\n                CVV:      "***",\n            }\n        }),\n    )\n\n    card := CreditCard{\n        Number:   "4111-1111-1111-1111",\n        ExpMonth: 12,\n        ExpYear:  2025,\n        CVV:      "123",\n    }\n\n    // Process the data\n    masked := c.Process(card)\n    fmt.Printf("%+v\\n", masked)\n    // TEXT output: {Number:****-****-****-1111 ExpMonth:12 ExpYear:2025 CVV:***}\n    // JSON output: {"Number":"****-****-****-1111","ExpMonth":12,"ExpYear":2025,"CVV":"***"}\n}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"circular-references",children:"Circular References"}),"\n",(0,s.jsx)(n.p,{children:"Censor can handle circular references in your data structures:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'type Node struct {\n    ID       string `censor:"display"`\n    Name     string\n    Parent   *Node\n    Children []*Node\n}\n\n// Create nodes with circular references\nparent := &Node{\n    ID:       "1",\n    Name:     "Parent",\n    Children: make([]*Node, 0),\n}\n\nchild := &Node{\n    ID:     "2",\n    Name:   "Child",\n    Parent: parent,\n}\n\nparent.Children = append(parent.Children, child)\n\n// Censor processes this correctly, preventing infinite recursion\n// TEXT output: {ID:1 Name:[CENSORED] Parent:<nil> Children:[{ID:2 Name:[CENSORED] Parent:0xc00010e100 Children:[]}]}\n// JSON output: {"ID":"1","Name":"[CENSORED]","Parent":null,"Children":[{"ID":"2","Name":"[CENSORED]","Parent":{"ID":"1","Name":"[CENSORED]","Parent":null,"Children":null},"Children":[]}]}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"reflection-edge-cases",children:"Reflection Edge Cases"}),"\n",(0,s.jsx)(n.p,{children:"Censor handles various reflection edge cases:"}),"\n",(0,s.jsx)(n.h3,{id:"zero-values",children:"Zero Values"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'// Zero values\nvar s string    // ""\nvar i int       // 0\nvar b bool      // false\nvar f float64   // 0.0\nvar t time.Time // Zero time\n\n// String zero value will not be masked\n// TEXT output: \n// JSON output: ""\n\n// Other zero values are displayed as is\n// TEXT output for i: 0\n// TEXT output for b: false\n// TEXT output for f: 0\n// TEXT output for t: 0001-01-01T00:00:00Z\n'})}),"\n",(0,s.jsx)(n.h3,{id:"unexported-fields",children:"Unexported Fields"}),"\n",(0,s.jsx)(n.p,{children:"Censor respects Go's visibility rules and doesn't access unexported fields:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'type User struct {\n    ID       string `censor:"display"`\n    Email    string\n    password string // unexported\n}\n\nuser := User{\n    ID:       "123",\n    Email:    "user@example.com",\n    password: "secret123",\n}\n\n// The unexported password field is not included in the output\n// TEXT output: {ID:123 Email:[CENSORED]}\n// JSON output: {"ID":"123","Email":"[CENSORED]"}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"complete-example",children:"Complete Example"}),"\n",(0,s.jsx)(n.p,{children:"Here's a complete example showing how Censor handles various special types:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'package main\n\nimport (\n    "fmt"\n    "github.com/vpakhuchyi/censor"\n)\n\n// Type alias\ntype Email string\n\n// Custom type\ntype CreditCard struct {\n    Number   string\n    ExpMonth int\n    ExpYear  int\n    CVV      string\n}\n\n// User with various special types\ntype User struct {\n    ID        string     `censor:"display"`\n    Email     Email\n    Card      CreditCard\n    Data      interface{}\n    AnyData   any\n    Temporary interface{} // Will be nil\n}\n\nfunc main() {\n    // Create a Censor instance with custom type handler\n    c := censor.New(\n        censor.WithTypeHandler(func(cc CreditCard) interface{} {\n            // Mask all except last 4 digits\n            if len(cc.Number) > 4 {\n                return CreditCard{\n                    Number:   "****-****-****-" + cc.Number[len(cc.Number)-4:],\n                    ExpMonth: cc.ExpMonth,\n                    ExpYear:  cc.ExpYear,\n                    CVV:      "***",\n                }\n            }\n            return CreditCard{\n                Number:   "[CENSORED]",\n                ExpMonth: cc.ExpMonth,\n                ExpYear:  cc.ExpYear,\n                CVV:      "***",\n            }\n        }),\n    )\n\n    // Create a user\n    user := User{\n        ID:    "123",\n        Email: "user@example.com",\n        Card: CreditCard{\n            Number:   "4111-1111-1111-1111",\n            ExpMonth: 12,\n            ExpYear:  2025,\n            CVV:      "123",\n        },\n        Data: map[string]string{\n            "api_key": "sk_live_123456789",\n        },\n        AnyData: "sensitive-data",\n    }\n\n    // Process the data\n    masked := c.Process(user)\n    fmt.Printf("%+v\\n", masked)\n    // TEXT output: {ID:123 Email:[CENSORED] Card:{Number:****-****-****-1111 ExpMonth:12 ExpYear:2025 CVV:***} Data:map[api_key:[CENSORED]] AnyData:[CENSORED] Temporary:<nil>}\n    // JSON output: {"ID":"123","Email":"[CENSORED]","Card":{"Number":"****-****-****-1111","ExpMonth":12,"ExpYear":2025,"CVV":"***"},"Data":{"api_key":"[CENSORED]"},"AnyData":"[CENSORED]","Temporary":null}\n}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"next-steps",children:"Next Steps"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsxs)(n.li,{children:["Learn about ",(0,s.jsx)(n.a,{href:"basic-types",children:"Basic Types"})]}),"\n",(0,s.jsxs)(n.li,{children:["Check out ",(0,s.jsx)(n.a,{href:"complex-types",children:"Complex Types"})]}),"\n",(0,s.jsxs)(n.li,{children:["See ",(0,s.jsx)(n.a,{href:"format-specific",children:"Format-Specific"})," handling"]}),"\n"]})]})}function d(e={}){const{wrapper:n}={...(0,a.R)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(p,{...e})}):p(e)}}}]);