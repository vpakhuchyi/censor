"use strict";(self.webpackChunkdocumentation=self.webpackChunkdocumentation||[]).push([[579],{6640:(n,e,s)=>{s.r(e),s.d(e,{assets:()=>c,contentTitle:()=>a,default:()=>g,frontMatter:()=>i,metadata:()=>r,toc:()=>d});const r=JSON.parse('{"id":"output-formats","title":"Output Formats","description":"Censor supports two primary output formats: Text and JSON. Each format is designed for specific use cases and provides different formatting options.","source":"@site/docs/output-formats.md","sourceDirName":".","slug":"/output-formats","permalink":"/censor/output-formats","draft":false,"unlisted":false,"editUrl":"https://github.com/vpakhuchyi/censor/tree/main/documentation/docs/output-formats.md","tags":[],"version":"current","lastUpdatedBy":"Viktor","lastUpdatedAt":1743858422000,"sidebarPosition":5,"frontMatter":{"sidebar_position":5}}');var t=s(4848),o=s(8453);const i={sidebar_position:5},a="Output Formats",c={},d=[{value:"Overview",id:"overview",level:2},{value:"Text Format",id:"text-format",level:3},{value:"JSON Format",id:"json-format",level:3},{value:"Text Format",id:"text-format-1",level:2},{value:"Basic Usage",id:"basic-usage",level:3},{value:"Text Format Features",id:"text-format-features",level:3},{value:"JSON Format",id:"json-format-1",level:2},{value:"Basic Usage",id:"basic-usage-1",level:3},{value:"JSON Format Features",id:"json-format-features",level:3},{value:"Choosing the Right Format",id:"choosing-the-right-format",level:2},{value:"Use Text Format When:",id:"use-text-format-when",level:3},{value:"Use JSON Format When:",id:"use-json-format-when",level:3},{value:"Format-Specific Configuration",id:"format-specific-configuration",level:2},{value:"Text Format Configuration",id:"text-format-configuration",level:3},{value:"JSON Format Configuration",id:"json-format-configuration",level:3},{value:"Best Practices",id:"best-practices",level:2},{value:"Common Pitfalls",id:"common-pitfalls",level:2}];function l(n){const e={code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,o.R)(),...n.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(e.header,{children:(0,t.jsx)(e.h1,{id:"output-formats",children:"Output Formats"})}),"\n",(0,t.jsx)(e.p,{children:"Censor supports two primary output formats: Text and JSON. Each format is designed for specific use cases and provides different formatting options."}),"\n",(0,t.jsx)(e.h2,{id:"overview",children:"Overview"}),"\n",(0,t.jsx)(e.h3,{id:"text-format",children:"Text Format"}),"\n",(0,t.jsxs)(e.ul,{children:["\n",(0,t.jsx)(e.li,{children:"Human-readable output"}),"\n",(0,t.jsx)(e.li,{children:"Default format for direct printing and logging"}),"\n",(0,t.jsx)(e.li,{children:"Ideal for development and debugging"}),"\n",(0,t.jsx)(e.li,{children:"Supports custom formatting patterns"}),"\n"]}),"\n",(0,t.jsx)(e.h3,{id:"json-format",children:"JSON Format"}),"\n",(0,t.jsxs)(e.ul,{children:["\n",(0,t.jsx)(e.li,{children:"Machine-readable output"}),"\n",(0,t.jsx)(e.li,{children:"Suitable for structured logging and APIs"}),"\n",(0,t.jsx)(e.li,{children:"Preserves data types"}),"\n",(0,t.jsx)(e.li,{children:"Compatible with JSON parsers"}),"\n"]}),"\n",(0,t.jsx)(e.h2,{id:"text-format-1",children:"Text Format"}),"\n",(0,t.jsx)(e.p,{children:"The text format is designed for human readability and is the default output format."}),"\n",(0,t.jsx)(e.h3,{id:"basic-usage",children:"Basic Usage"}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'package main\n\nimport (\n    "fmt"\n    "github.com/vpakhuchyi/censor"\n)\n\nfunc main() {\n    c := censor.New()\n\n    type User struct {\n        ID       string `censor:"display"`\n        Email    string\n        Password string\n    }\n\n    user := User{\n        ID:       "123",\n        Email:    "user@example.com",\n        Password: "secret",\n    }\n\n    // Text format output\n    masked := c.Process(user)\n    fmt.Printf("%v\\n", masked)\n    // Output: {ID: 123, Email: [CENSORED], Password: [CENSORED]}\n}\n'})}),"\n",(0,t.jsx)(e.h3,{id:"text-format-features",children:"Text Format Features"}),"\n",(0,t.jsxs)(e.ol,{children:["\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Struct Formatting"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'type Address struct {\n    Street string\n    City   string\n}\n\ntype User struct {\n    ID      string  `censor:"display"`\n    Address Address\n}\n\nuser := User{\n    ID: "123",\n    Address: Address{\n        Street: "123 Main St",\n        City:   "Anytown",\n    },\n}\n\nfmt.Printf("%v\\n", c.Process(user))\n// Output: {ID: 123, Address: {Street: [CENSORED], City: [CENSORED]}}\n'})}),"\n"]}),"\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Slice and Array Formatting"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'data := []string{"secret1", "secret2"}\nfmt.Printf("%v\\n", c.Process(data))\n// Output: [[CENSORED], [CENSORED]]\n'})}),"\n"]}),"\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Map Formatting"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'data := map[string]string{\n    "visible": "public",\n    "hidden":  "secret",\n}\nfmt.Printf("%v\\n", c.Process(data))\n// Output: map[visible: public, hidden: [CENSORED]]\n'})}),"\n"]}),"\n"]}),"\n",(0,t.jsx)(e.h2,{id:"json-format-1",children:"JSON Format"}),"\n",(0,t.jsx)(e.p,{children:"The JSON format is designed for machine readability and structured data interchange."}),"\n",(0,t.jsx)(e.h3,{id:"basic-usage-1",children:"Basic Usage"}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'package main\n\nimport (\n    "encoding/json"\n    "fmt"\n    "github.com/vpakhuchyi/censor"\n)\n\nfunc main() {\n    c := censor.New()\n\n    type User struct {\n        ID       string `json:"id" censor:"display"`\n        Email    string `json:"email"`\n        Password string `json:"password"`\n    }\n\n    user := User{\n        ID:       "123",\n        Email:    "user@example.com",\n        Password: "secret",\n    }\n\n    // Process and convert to JSON\n    masked := c.Process(user)\n    jsonData, _ := json.Marshal(masked)\n    fmt.Println(string(jsonData))\n    // Output: {"id":"123","email":"[CENSORED]","password":"[CENSORED]"}\n}\n'})}),"\n",(0,t.jsx)(e.h3,{id:"json-format-features",children:"JSON Format Features"}),"\n",(0,t.jsxs)(e.ol,{children:["\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Struct to JSON"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'type Address struct {\n    Street string `json:"street"`\n    City   string `json:"city"`\n}\n\ntype User struct {\n    ID      string  `json:"id" censor:"display"`\n    Address Address `json:"address"`\n}\n\nuser := User{\n    ID: "123",\n    Address: Address{\n        Street: "123 Main St",\n        City:   "Anytown",\n    },\n}\n\nmasked := c.Process(user)\njsonData, _ := json.MarshalIndent(masked, "", "  ")\nfmt.Println(string(jsonData))\n// Output:\n// {\n//   "id": "123",\n//   "address": {\n//     "street": "[CENSORED]",\n//     "city": "[CENSORED]"\n//   }\n// }\n'})}),"\n"]}),"\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Arrays and Slices in JSON"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'type User struct {\n    ID     string   `json:"id" censor:"display"`\n    Emails []string `json:"emails"`\n}\n\nuser := User{\n    ID:     "123",\n    Emails: []string{"user1@example.com", "user2@example.com"},\n}\n\nmasked := c.Process(user)\njsonData, _ := json.Marshal(masked)\nfmt.Println(string(jsonData))\n// Output: {"id":"123","emails":["[CENSORED]","[CENSORED]"]}\n'})}),"\n"]}),"\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Maps in JSON"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'data := map[string]interface{}{\n    "id": "123",\n    "credentials": map[string]string{\n        "username": "john_doe",\n        "password": "secret",\n    },\n}\n\nmasked := c.Process(data)\njsonData, _ := json.MarshalIndent(masked, "", "  ")\nfmt.Println(string(jsonData))\n// Output:\n// {\n//   "id": "123",\n//   "credentials": {\n//     "username": "john_doe",\n//     "password": "[CENSORED]"\n//   }\n// }\n'})}),"\n"]}),"\n"]}),"\n",(0,t.jsx)(e.h2,{id:"choosing-the-right-format",children:"Choosing the Right Format"}),"\n",(0,t.jsx)(e.h3,{id:"use-text-format-when",children:"Use Text Format When:"}),"\n",(0,t.jsxs)(e.ul,{children:["\n",(0,t.jsx)(e.li,{children:"Debugging application behavior"}),"\n",(0,t.jsx)(e.li,{children:"Reading logs directly"}),"\n",(0,t.jsx)(e.li,{children:"Printing to console"}),"\n",(0,t.jsx)(e.li,{children:"Need human-readable output"}),"\n",(0,t.jsx)(e.li,{children:"Working with standard Go printing functions"}),"\n"]}),"\n",(0,t.jsx)(e.h3,{id:"use-json-format-when",children:"Use JSON Format When:"}),"\n",(0,t.jsxs)(e.ul,{children:["\n",(0,t.jsx)(e.li,{children:"Integrating with logging systems"}),"\n",(0,t.jsx)(e.li,{children:"Building APIs"}),"\n",(0,t.jsx)(e.li,{children:"Storing structured data"}),"\n",(0,t.jsx)(e.li,{children:"Need machine-readable output"}),"\n",(0,t.jsx)(e.li,{children:"Working with JSON-based tools and services"}),"\n"]}),"\n",(0,t.jsx)(e.h2,{id:"format-specific-configuration",children:"Format-Specific Configuration"}),"\n",(0,t.jsx)(e.h3,{id:"text-format-configuration",children:"Text Format Configuration"}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'cfg := censor.Config{\n    Formatter: censor.FormatterConfig{\n        MaskValue: "***",            // Custom mask for text output\n        ExcludePatterns: []string{\n            `\\d{4}-\\d{4}-\\d{4}-\\d{4}`,\n        },\n    },\n}\n\nc := censor.NewWithOpts(censor.WithConfig(&cfg))\n'})}),"\n",(0,t.jsx)(e.h3,{id:"json-format-configuration",children:"JSON Format Configuration"}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'cfg := censor.Config{\n    Formatter: censor.FormatterConfig{\n        MaskValue: "\\"[MASKED]\\"",  // Custom mask for JSON output\n        ExcludePatterns: []string{\n            `\\d{4}-\\d{4}-\\d{4}-\\d{4}`,\n        },\n    },\n}\n\nc := censor.NewWithOpts(censor.WithConfig(&cfg))\n'})}),"\n",(0,t.jsx)(e.h2,{id:"best-practices",children:"Best Practices"}),"\n",(0,t.jsxs)(e.ol,{children:["\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Consistent Format Usage"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'// Good: Consistent text format\nlogger.Printf("User: %v", c.Process(user))\nlogger.Printf("Payment: %v", c.Process(payment))\n\n// Good: Consistent JSON format\njsonLogger.Log(c.Process(user))\njsonLogger.Log(c.Process(payment))\n\n// Bad: Mixed formats\nlogger.Printf("User: %v", c.Process(user))\njsonLogger.Log(c.Process(user))\n'})}),"\n"]}),"\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Format-Specific Masking"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'// Text format configuration\ntextCfg := censor.Config{\n    Formatter: censor.FormatterConfig{\n        MaskValue: "***",\n    },\n}\n\n// JSON format configuration\njsonCfg := censor.Config{\n    Formatter: censor.FormatterConfig{\n        MaskValue: "\\"[MASKED]\\"",\n    },\n}\n'})}),"\n"]}),"\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Logger Integration"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'// Text format for standard logging\nstdLogger := log.New(os.Stdout, "", log.LstdFlags)\nstdLogger.Printf("Data: %v", c.Process(data))\n\n// JSON format for structured logging\njsonLogger := zerolog.New(os.Stdout)\njsonLogger.Info().Interface("data", c.Process(data)).Send()\n'})}),"\n"]}),"\n"]}),"\n",(0,t.jsx)(e.h2,{id:"common-pitfalls",children:"Common Pitfalls"}),"\n",(0,t.jsxs)(e.ol,{children:["\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Mixing Formats"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'// Problem: Inconsistent output format\ndata := map[string]string{"password": "secret"}\n\nfmt.Printf("%v\\n", c.Process(data))        // Text format\njson.NewEncoder(w).Encode(c.Process(data)) // JSON format\n\n// Solution: Use consistent format for each context\ntextLogger := log.New(os.Stdout, "", 0)\njsonLogger := zerolog.New(os.Stdout)\n\ntextLogger.Printf("%v", c.Process(data))   // Always text\njsonLogger.Info().Interface("data", c.Process(data)).Send() // Always JSON\n'})}),"\n"]}),"\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Format-Specific Escaping"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'// Problem: Incorrect escaping\ncfg := censor.Config{\n    Formatter: censor.FormatterConfig{\n        MaskValue: "[CENSORED]", // Might cause JSON parsing issues\n    },\n}\n\n// Solution: Use format-specific masks\ntextCfg := censor.Config{\n    Formatter: censor.FormatterConfig{\n        MaskValue: "[CENSORED]",\n    },\n}\n\njsonCfg := censor.Config{\n    Formatter: censor.FormatterConfig{\n        MaskValue: "\\"[CENSORED]\\"",\n    },\n}\n'})}),"\n"]}),"\n",(0,t.jsxs)(e.li,{children:["\n",(0,t.jsx)(e.p,{children:(0,t.jsx)(e.strong,{children:"Type Preservation"})}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'// Problem: Type mismatch in JSON\ntype User struct {\n    Age int `json:"age"`\n}\n\n// Solution: Use appropriate types for each format\ntextAge := fmt.Sprintf("%v", c.Process(user.Age))\njsonAge := c.Process(user.Age).(int)\n'})}),"\n"]}),"\n"]})]})}function g(n={}){const{wrapper:e}={...(0,o.R)(),...n.components};return e?(0,t.jsx)(e,{...n,children:(0,t.jsx)(l,{...n})}):l(n)}},8453:(n,e,s)=>{s.d(e,{R:()=>i,x:()=>a});var r=s(6540);const t={},o=r.createContext(t);function i(n){const e=r.useContext(o);return r.useMemo((function(){return"function"==typeof n?n(e):{...e,...n}}),[e,n])}function a(n){let e;return e=n.disableParentContext?"function"==typeof n.components?n.components(t):n.components||t:i(n.components),r.createElement(o.Provider,{value:e},n.children)}}}]);