"use strict";(self.webpackChunkdocumentation=self.webpackChunkdocumentation||[]).push([[417],{8379:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>c,contentTitle:()=>o,default:()=>p,frontMatter:()=>a,metadata:()=>s,toc:()=>l});const s=JSON.parse('{"id":"type-handling/format-specific","title":"Format-Specific Type Handling","description":"Censor supports different output formats (TEXT and JSON), which can sometimes handle the same Go types differently. This guide explains the format-specific behavior of Censor.","source":"@site/docs/type-handling/format-specific.md","sourceDirName":"type-handling","slug":"/type-handling/format-specific","permalink":"/censor/type-handling/format-specific","draft":false,"unlisted":false,"editUrl":"https://github.com/vpakhuchyi/censor-doc/tree/main/docs/type-handling/format-specific.md","tags":[],"version":"current","lastUpdatedBy":"Viktor","lastUpdatedAt":1743858422000,"sidebarPosition":4,"frontMatter":{"sidebar_position":4},"sidebar":"docs","previous":{"title":"Special Types","permalink":"/censor/type-handling/special-types"},"next":{"title":"Zap","permalink":"/censor/logger-handlers/zap"}}');var r=t(4848),i=t(8453);const a={sidebar_position:4},o="Format-Specific Type Handling",c={},l=[{value:"TEXT vs JSON Output",id:"text-vs-json-output",level:2},{value:"Type-Specific Behavior",id:"type-specific-behavior",level:2},{value:"Strings",id:"strings",level:3},{value:"Maps",id:"maps",level:3},{value:"Structs",id:"structs",level:3},{value:"Time",id:"time",level:3},{value:"Slices and Arrays",id:"slices-and-arrays",level:3},{value:"Pointers",id:"pointers",level:3},{value:"<code>nil</code> Values",id:"nil-values",level:3},{value:"JSON Struct Tags",id:"json-struct-tags",level:2},{value:"Changing Format at Runtime",id:"changing-format-at-runtime",level:2},{value:"Complete Example",id:"complete-example",level:2},{value:"Next Steps",id:"next-steps",level:2}];function d(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",...(0,i.R)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(n.header,{children:(0,r.jsx)(n.h1,{id:"format-specific-type-handling",children:"Format-Specific Type Handling"})}),"\n",(0,r.jsx)(n.p,{children:"Censor supports different output formats (TEXT and JSON), which can sometimes handle the same Go types differently. This guide explains the format-specific behavior of Censor."}),"\n",(0,r.jsx)(n.h2,{id:"text-vs-json-output",children:"TEXT vs JSON Output"}),"\n",(0,r.jsx)(n.p,{children:"The two main formats supported by Censor are:"}),"\n",(0,r.jsxs)(n.ol,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"TEXT"}),": Human-readable format, similar to how Go prints values with ",(0,r.jsx)(n.code,{children:"fmt.Print()"})]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"JSON"}),": Structured format suitable for APIs, logs, and data interchange"]}),"\n"]}),"\n",(0,r.jsx)(n.h2,{id:"type-specific-behavior",children:"Type-Specific Behavior"}),"\n",(0,r.jsx)(n.h3,{id:"strings",children:"Strings"}),"\n",(0,r.jsx)(n.p,{children:"Strings are processed similarly in both formats, but with different outputs:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'// Create a Censor instance with TEXT format\ncText := censor.New(censor.WithFormatter(censor.FormatterText))\n\n// Create a Censor instance with JSON format\ncJSON := censor.New(censor.WithFormatter(censor.FormatterJSON))\n\n// Process a string\nemail := "user@example.com"\n\n// TEXT output: [CENSORED]\n// JSON output: "[CENSORED]"  // Note the quotes for JSON strings\n'})}),"\n",(0,r.jsx)(n.h3,{id:"maps",children:"Maps"}),"\n",(0,r.jsx)(n.p,{children:"Maps are represented differently in TEXT and JSON:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'// Map with string keys and values\nuser := map[string]string{\n    "id":       "123",\n    "email":    "user@example.com",\n    "password": "secret123",\n}\n\n// TEXT output: map[id:123 email:[CENSORED] password:[CENSORED]]\n// JSON output: {"id":"123","email":"[CENSORED]","password":"[CENSORED]"}\n'})}),"\n",(0,r.jsx)(n.p,{children:"Maps with non-string keys are handled differently:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'// Map with integer keys\nscores := map[int]int{\n    1: 100,\n    2: 200,\n    3: 300,\n}\n\n// TEXT output: map[1:100 2:200 3:300]\n// JSON output: {"1":100,"2":200,"3":300}  // Keys are converted to strings in JSON\n'})}),"\n",(0,r.jsx)(n.h3,{id:"structs",children:"Structs"}),"\n",(0,r.jsx)(n.p,{children:"Structs are displayed differently in TEXT and JSON:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'type User struct {\n    ID        string `json:"id" censor:"display"`\n    Email     string `json:"email"`\n    Password  string `json:"password"`\n}\n\nuser := User{\n    ID:       "123",\n    Email:    "user@example.com",\n    Password: "secret123",\n}\n\n// TEXT output: {ID:123 Email:[CENSORED] Password:[CENSORED]}\n// JSON output: {"id":"123","email":"[CENSORED]","password":"[CENSORED]"}\n'})}),"\n",(0,r.jsxs)(n.p,{children:["Note how JSON respects the ",(0,r.jsx)(n.code,{children:"json"})," struct tags, while TEXT uses the field names."]}),"\n",(0,r.jsx)(n.h3,{id:"time",children:"Time"}),"\n",(0,r.jsx)(n.p,{children:"Time values are formatted according to the configuration:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)\n\n// Default format\n// TEXT output: 2023-01-01T12:00:00Z\n// JSON output: "2023-01-01T12:00:00Z"  // JSON wraps the time string in quotes\n\n// With custom time format\nc := censor.New(censor.WithTimeFormat("2006-01-02"))\n// TEXT output: 2023-01-01\n// JSON output: "2023-01-01"\n'})}),"\n",(0,r.jsx)(n.h3,{id:"slices-and-arrays",children:"Slices and Arrays"}),"\n",(0,r.jsx)(n.p,{children:"Slices and arrays are formatted differently:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'// Slice of strings\nemails := []string{\n    "user1@example.com",\n    "user2@example.com",\n    "user3@example.com",\n}\n\n// TEXT output: [[CENSORED] [CENSORED] [CENSORED]]\n// JSON output: ["[CENSORED]","[CENSORED]","[CENSORED]"]\n\n// Slice of integers\nscores := []int{100, 200, 300}\n\n// TEXT output: [100 200 300]\n// JSON output: [100,200,300]\n'})}),"\n",(0,r.jsx)(n.h3,{id:"pointers",children:"Pointers"}),"\n",(0,r.jsx)(n.p,{children:"Pointers are dereferenced in both formats:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'email := "user@example.com"\nemailPtr := &email\n\n// TEXT output: [CENSORED]\n// JSON output: "[CENSORED]"\n\n// Nil pointers\nvar nilPtr *string\n\n// TEXT output: <nil>\n// JSON output: null\n'})}),"\n",(0,r.jsxs)(n.h3,{id:"nil-values",children:[(0,r.jsx)(n.code,{children:"nil"})," Values"]}),"\n",(0,r.jsx)(n.p,{children:"Nil values are handled appropriately for each format:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:"// Nil interface\nvar data interface{}\n\n// TEXT output: <nil>\n// JSON output: null\n\n// Nil slice\nvar slice []string\n\n// TEXT output: []\n// JSON output: []\n"})}),"\n",(0,r.jsx)(n.h2,{id:"json-struct-tags",children:"JSON Struct Tags"}),"\n",(0,r.jsxs)(n.p,{children:["For JSON output, Censor respects the ",(0,r.jsx)(n.code,{children:"json"})," struct tags:"]}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'type User struct {\n    ID           string `json:"id" censor:"display"`\n    Email        string `json:"email"`\n    Password     string `json:"password"`\n    InternalNote string `json:"-"`  // Will be omitted in JSON output\n    APIKey       string `json:"api_key,omitempty"`  // Will be omitted in JSON if empty\n}\n\nuser := User{\n    ID:       "123",\n    Email:    "user@example.com",\n    Password: "secret123",\n    InternalNote: "This is an internal note",\n    // APIKey is empty, so it will be omitted with omitempty\n}\n\n// TEXT output: {ID:123 Email:[CENSORED] Password:[CENSORED] InternalNote:[CENSORED] APIKey:}\n// JSON output: {"id":"123","email":"[CENSORED]","password":"[CENSORED]"}\n'})}),"\n",(0,r.jsx)(n.h2,{id:"changing-format-at-runtime",children:"Changing Format at Runtime"}),"\n",(0,r.jsx)(n.p,{children:"You can change the output format at runtime:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'package main\n\nimport (\n    "fmt"\n    "github.com/vpakhuchyi/censor"\n)\n\nfunc main() {\n    // Create a Censor instance with TEXT format\n    c := censor.New(censor.WithFormatter(censor.FormatterText))\n\n    type User struct {\n        ID       string `json:"id" censor:"display"`\n        Email    string `json:"email"`\n        Password string `json:"password"`\n    }\n\n    user := User{\n        ID:       "123",\n        Email:    "user@example.com",\n        Password: "secret123",\n    }\n\n    // Process with TEXT format\n    textOutput := c.Process(user)\n    fmt.Println("TEXT output:", textOutput)\n    // TEXT output: {ID:123 Email:[CENSORED] Password:[CENSORED]}\n\n    // Change to JSON format\n    c = censor.New(censor.WithFormatter(censor.FormatterJSON))\n    \n    // Process with JSON format\n    jsonOutput := c.Process(user)\n    fmt.Println("JSON output:", jsonOutput)\n    // JSON output: {"id":"123","email":"[CENSORED]","password":"[CENSORED]"}\n}\n'})}),"\n",(0,r.jsx)(n.h2,{id:"complete-example",children:"Complete Example"}),"\n",(0,r.jsx)(n.p,{children:"Here's a complete example showing format-specific behavior:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'package main\n\nimport (\n    "fmt"\n    "time"\n    "github.com/vpakhuchyi/censor"\n)\n\ntype Address struct {\n    Street  string `json:"street"`\n    City    string `json:"city"`\n    ZipCode string `json:"zip_code"`\n}\n\ntype User struct {\n    ID        string    `json:"id" censor:"display"`\n    Email     string    `json:"email"`\n    Password  string    `json:"password"`\n    CreatedAt time.Time `json:"created_at"`\n    Address   Address   `json:"address"`\n    Tags      []string  `json:"tags"`\n    Metadata  map[string]string `json:"metadata"`\n}\n\nfunc main() {\n    // Create users\n    user := User{\n        ID:        "123",\n        Email:     "user@example.com",\n        Password:  "secret123",\n        CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),\n        Address: Address{\n            Street:  "123 Main St",\n            City:    "New York",\n            ZipCode: "10001",\n        },\n        Tags: []string{"premium", "active"},\n        Metadata: map[string]string{\n            "last_login": "2023-01-01",\n            "api_key":    "sk_live_123456789",\n        },\n    }\n\n    // TEXT format\n    cText := censor.New(censor.WithFormatter(censor.FormatterText))\n    textOutput := cText.Process(user)\n    fmt.Println("TEXT output:", textOutput)\n    // TEXT output: {ID:123 Email:[CENSORED] Password:[CENSORED] CreatedAt:2023-01-01 12:00:00 +0000 UTC Address:{Street:[CENSORED] City:[CENSORED] ZipCode:[CENSORED]} Tags:[[CENSORED] [CENSORED]] Metadata:map[api_key:[CENSORED] last_login:[CENSORED]]}\n\n    // JSON format\n    cJSON := censor.New(censor.WithFormatter(censor.FormatterJSON))\n    jsonOutput := cJSON.Process(user)\n    fmt.Println("JSON output:", jsonOutput)\n    // JSON output: {"id":"123","email":"[CENSORED]","password":"[CENSORED]","created_at":"2023-01-01T12:00:00Z","address":{"street":"[CENSORED]","city":"[CENSORED]","zip_code":"[CENSORED]"},"tags":["[CENSORED]","[CENSORED]"],"metadata":{"api_key":"[CENSORED]","last_login":"[CENSORED]"}}\n}\n'})}),"\n",(0,r.jsx)(n.h2,{id:"next-steps",children:"Next Steps"}),"\n",(0,r.jsxs)(n.ol,{children:["\n",(0,r.jsxs)(n.li,{children:["Learn about ",(0,r.jsx)(n.a,{href:"basic-types",children:"Basic Types"})]}),"\n",(0,r.jsxs)(n.li,{children:["Check out ",(0,r.jsx)(n.a,{href:"complex-types",children:"Complex Types"})]}),"\n",(0,r.jsxs)(n.li,{children:["See ",(0,r.jsx)(n.a,{href:"special-types",children:"Special Types"})]}),"\n"]})]})}function p(e={}){const{wrapper:n}={...(0,i.R)(),...e.components};return n?(0,r.jsx)(n,{...e,children:(0,r.jsx)(d,{...e})}):d(e)}},8453:(e,n,t)=>{t.d(n,{R:()=>a,x:()=>o});var s=t(6540);const r={},i=s.createContext(r);function a(e){const n=s.useContext(i);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:a(e.components),s.createElement(i.Provider,{value:n},e.children)}}}]);