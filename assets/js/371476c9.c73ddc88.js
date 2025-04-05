"use strict";(self.webpackChunkdocumentation=self.webpackChunkdocumentation||[]).push([[697],{1818:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>t,contentTitle:()=>l,default:()=>g,frontMatter:()=>i,metadata:()=>o,toc:()=>c});const o=JSON.parse('{"id":"logger-handlers/zap","title":"Zap Logger Integration","description":"This guide covers how to integrate Censor with the Zap logging framework.","source":"@site/docs/logger-handlers/zap.md","sourceDirName":"logger-handlers","slug":"/logger-handlers/zap","permalink":"/censor/logger-handlers/zap","draft":false,"unlisted":false,"editUrl":"https://github.com/vpakhuchyi/censor/tree/main/documentation/docs/logger-handlers/zap.md","tags":[],"version":"current","lastUpdatedBy":"Viktor","lastUpdatedAt":1743858422000,"sidebarPosition":2,"frontMatter":{"sidebar_position":2},"sidebar":"docs","previous":{"title":"Format-Specific","permalink":"/censor/type-handling/format-specific"},"next":{"title":"Slog","permalink":"/censor/logger-handlers/slog"}}');var s=r(4848),a=r(8453);const i={sidebar_position:2},l="Zap Logger Integration",t={},c=[{value:"Installation",id:"installation",level:2},{value:"Basic Setup",id:"basic-setup",level:2},{value:"Configuration Options",id:"configuration-options",level:2},{value:"Custom Mask Value",id:"custom-mask-value",level:3},{value:"Exclude Patterns",id:"exclude-patterns",level:3},{value:"JSON Output",id:"json-output",level:3},{value:"Advanced Usage",id:"advanced-usage",level:2},{value:"Custom Field Processors",id:"custom-field-processors",level:3},{value:"Field Name Mapping",id:"field-name-mapping",level:3},{value:"Complete Example",id:"complete-example",level:2},{value:"Best Practices",id:"best-practices",level:2},{value:"Next Steps",id:"next-steps",level:2}];function d(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,a.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.header,{children:(0,s.jsx)(n.h1,{id:"zap-logger-integration",children:"Zap Logger Integration"})}),"\n",(0,s.jsx)(n.p,{children:"This guide covers how to integrate Censor with the Zap logging framework."}),"\n",(0,s.jsx)(n.h2,{id:"installation",children:"Installation"}),"\n",(0,s.jsx)(n.p,{children:"First, install the required dependencies:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"go get -u go.uber.org/zap\ngo get -u github.com/vpakhuchyi/censor\n"})}),"\n",(0,s.jsx)(n.h2,{id:"basic-setup",children:"Basic Setup"}),"\n",(0,s.jsx)(n.p,{children:"Here's a basic example of using Censor with Zap:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'package main\n\nimport (\n    "os"\n    "go.uber.org/zap"\n    "go.uber.org/zap/zapcore"\n    "github.com/vpakhuchyi/censor"\n    "github.com/vpakhuchyi/censor/logger/zap"\n)\n\nfunc main() {\n    // Create a Censor instance\n    c := censor.New()\n\n    // Create a Zap logger with Censor handler\n    logger := zap.New(\n        censorlog.NewHandler(\n            zapcore.NewCore(\n                zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),\n                zapcore.AddSync(os.Stdout),\n                zapcore.DebugLevel,\n            ),\n        ),\n    )\n\n    // Log sensitive data\n    user := User{\n        ID:       "123",\n        Email:    "user@example.com",\n        Password: "secret123",\n    }\n\n    logger.Info("user data", zap.Any("user", user))\n}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"configuration-options",children:"Configuration Options"}),"\n",(0,s.jsx)(n.h3,{id:"custom-mask-value",children:"Custom Mask Value"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'logger := zap.New(\n    censorlog.NewHandler(\n        zapcore.NewCore(...),\n        censorlog.WithMaskValue("[REDACTED]"),\n    ),\n)\n'})}),"\n",(0,s.jsx)(n.h3,{id:"exclude-patterns",children:"Exclude Patterns"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:"logger := zap.New(\n    censorlog.NewHandler(\n        zapcore.NewCore(...),\n        censorlog.WithExcludePatterns([]string{\n            `\\d{4}-\\d{4}-\\d{4}-\\d{4}`, // Credit card numbers\n            `\\d{3}-\\d{2}-\\d{4}`,       // SSN\n        }),\n    ),\n)\n"})}),"\n",(0,s.jsx)(n.h3,{id:"json-output",children:"JSON Output"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:"logger := zap.New(\n    censorlog.NewHandler(\n        zapcore.NewCore(\n            zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),\n            zapcore.AddSync(os.Stdout),\n            zapcore.DebugLevel,\n        ),\n    ),\n)\n"})}),"\n",(0,s.jsx)(n.h2,{id:"advanced-usage",children:"Advanced Usage"}),"\n",(0,s.jsx)(n.h3,{id:"custom-field-processors",children:"Custom Field Processors"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'type CustomType string\n\nlogger := zap.New(\n    censorlog.NewHandler(\n        zapcore.NewCore(...),\n        censorlog.WithTypeHandler(reflect.TypeOf(CustomType("")), func(v interface{}) string {\n            return "[CUSTOM_MASKED]"\n        }),\n    ),\n)\n'})}),"\n",(0,s.jsx)(n.h3,{id:"field-name-mapping",children:"Field Name Mapping"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:"logger := zap.New(\n    censorlog.NewHandler(\n        zapcore.NewCore(...),\n        censorlog.WithFieldNameMapper(func(name string) string {\n            return strings.ToUpper(name)\n        }),\n    ),\n)\n"})}),"\n",(0,s.jsx)(n.h2,{id:"complete-example",children:"Complete Example"}),"\n",(0,s.jsx)(n.p,{children:"Here's a complete example showing various Zap integration features:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:'package main\n\nimport (\n    "os"\n    "go.uber.org/zap"\n    "go.uber.org/zap/zapcore"\n    "github.com/vpakhuchyi/censor"\n    "github.com/vpakhuchyi/censor/logger/zap"\n)\n\ntype User struct {\n    ID       string `censor:"display"`\n    Email    string\n    Password string `censor:"mask"`\n}\n\nfunc main() {\n    // Create a Censor instance with custom configuration\n    c := censor.New(\n        censor.WithMaskValue("[REDACTED]"),\n        censor.WithExcludePatterns([]string{\n            `\\d{4}-\\d{4}-\\d{4}-\\d{4}`,\n            `\\d{3}-\\d{2}-\\d{4}`,\n        }),\n    )\n\n    // Create encoder config\n    encoderConfig := zap.NewProductionEncoderConfig()\n    encoderConfig.TimeKey = "timestamp"\n    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder\n\n    // Create a Zap logger with Censor handler\n    logger := zap.New(\n        censorlog.NewHandler(\n            zapcore.NewCore(\n                zapcore.NewJSONEncoder(encoderConfig),\n                zapcore.AddSync(os.Stdout),\n                zapcore.DebugLevel,\n            ),\n            censorlog.WithMaskValue("[REDACTED]"),\n            censorlog.WithExcludePatterns([]string{\n                `\\d{4}-\\d{4}-\\d{4}-\\d{4}`,\n                `\\d{3}-\\d{2}-\\d{4}`,\n            }),\n        ),\n    )\n\n    // Create a user\n    user := User{\n        ID:       "123",\n        Email:    "user@example.com",\n        Password: "secret123",\n    }\n\n    // Log with different levels\n    logger.Debug("debug message", zap.Any("user", user))\n    logger.Info("info message", zap.Any("user", user))\n    logger.Warn("warn message", zap.Any("user", user))\n    logger.Error("error message", zap.Any("user", user))\n\n    // Log with fields\n    logger.Info("user action",\n        zap.String("action", "login"),\n        zap.Any("user", user),\n        zap.String("ip", "192.168.1.1"),\n    )\n\n    // Log with error\n    err := fmt.Errorf("database error")\n    logger.Error("operation failed",\n        zap.Error(err),\n        zap.Any("user", user),\n    )\n}\n'})}),"\n",(0,s.jsx)(n.h2,{id:"best-practices",children:"Best Practices"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsxs)(n.li,{children:["\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"Logger Initialization"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"Initialize logger once at application startup"}),"\n",(0,s.jsx)(n.li,{children:"Use appropriate log level for environment"}),"\n",(0,s.jsx)(n.li,{children:"Configure output format based on needs"}),"\n"]}),"\n"]}),"\n",(0,s.jsxs)(n.li,{children:["\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"Structured Logging"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"Use structured fields instead of string formatting"}),"\n",(0,s.jsx)(n.li,{children:"Include context with each log entry"}),"\n",(0,s.jsx)(n.li,{children:"Use appropriate field types"}),"\n"]}),"\n"]}),"\n",(0,s.jsxs)(n.li,{children:["\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"Error Handling"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"Log errors with context"}),"\n",(0,s.jsx)(n.li,{children:"Use appropriate log levels"}),"\n",(0,s.jsx)(n.li,{children:"Include stack traces when needed"}),"\n"]}),"\n"]}),"\n",(0,s.jsxs)(n.li,{children:["\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"Performance"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"Use appropriate encoder"}),"\n",(0,s.jsx)(n.li,{children:"Configure buffer sizes"}),"\n",(0,s.jsx)(n.li,{children:"Avoid unnecessary string formatting"}),"\n"]}),"\n"]}),"\n"]}),"\n",(0,s.jsx)(n.h2,{id:"next-steps",children:"Next Steps"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsxs)(n.li,{children:["Learn about ",(0,s.jsx)(n.a,{href:"slog",children:"Slog Integration"})]}),"\n",(0,s.jsxs)(n.li,{children:["Check out ",(0,s.jsx)(n.a,{href:"zerolog",children:"Zerolog Integration"})]}),"\n",(0,s.jsxs)(n.li,{children:["Review ",(0,s.jsx)(n.a,{href:"../configuration",children:"Configuration"})," options"]}),"\n",(0,s.jsxs)(n.li,{children:["See ",(0,s.jsx)(n.a,{href:"../examples/data-leak-prevention",children:"Examples"})," for more use cases"]}),"\n"]})]})}function g(e={}){const{wrapper:n}={...(0,a.R)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},8453:(e,n,r)=>{r.d(n,{R:()=>i,x:()=>l});var o=r(6540);const s={},a=o.createContext(s);function i(e){const n=o.useContext(a);return o.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function l(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:i(e.components),o.createElement(a.Provider,{value:n},e.children)}}}]);