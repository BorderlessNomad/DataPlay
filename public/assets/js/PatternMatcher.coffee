class window.PatternMatcher
  @entities: [
    {
      name: 'date'
      handlers: [
        {
          pattern: /^(((19|20)?\d\d)[- \/.])?((0)?([1-9])|1[012])[- \/.]((0)?([1-9])|[12][0-9]|3[01])$/
          parse: (m) => @parse m
        }
        {
          pattern: /^(0[1-9]|1[012])([- \/.])(0[1-9]|[12][0-9]|3[01])\2(19|20)\d\d$/
          parse: (src, s) ->
            vals = src.split s
            "#{vals[2]}-#{vals[0]}-#{vals[1]}" 
        }
        {
          pattern: /^(0[1-9]|[12][0-9]|3[01])([- \/.])(0[1-9]|1[012])\2(19|20)\d\d$/
          parse: (src, s) ->
            vals = src.split s
            "#{vals[2]}-#{vals[1]}-#{vals[0]}" 
        }
        {
          pattern: /^([1-9]|[12][0-9]|3[01])([- \/.])([1-9]|1[012])\2(19|20)\d\d$/
          parse: (src, s) ->
            vals = src.split s
            "#{vals[2]}-0#{vals[1]}-0#{vals[0]}" 
        }
        {
          pattern: /^([1-9]|[12][0-9]|3[01])([- \/.])(0[1-9]|1[012])\2(19|20)\d\d$/
          parse: (src, s) ->
            vals = src.split s
            "#{vals[2]}-#{vals[1]}-0#{vals[0]}" 
        }
        {
          pattern: /^(0[1-9]|[12][0-9]|3[01])([- \/.])([1-9]|1[012])\2(19|20)\d\d$/
          parse: (src, s) ->
            vals = src.split s
            "#{vals[2]}-0#{vals[1]}-#{vals[0]}" 
        }
        {
          pattern: /^(0[1-9]|[12][0-9]|3[01])([- \/.])(0[1-9]|1[012])\2\d\d$/
          parse: (src, s) ->
            vals = src.split s
            cy = new Date().getYear()
            y = parseInt("#20{vals[2]}")
            year = if cy>y then y else "#19{vals[2]}"
            "#{year}-#{vals[1]}-#{vals[0]}" 
        }
        {
          pattern: /^\d\d([- \/.])(0[1-9]|1[012])\2(0[1-9]|[12][0-9]|3[01])$/
          parse: (src, s) ->
            vals = src.split s
            cy = new Date().getYear()
            y = parseInt("#20{vals[0]}")
            year = if cy>y then y else "#19{vals[0]}"
            "#{year}-#{vals[1]}-#{vals[2]}" 
        }
        {
          pattern: /^(0[1-9]|1[012])([- \/.])(0[1-9]|[12][0-9]|3[01])\2\d\d$/
          parse: (src, s) ->
            vals = src.split s
            cy = new Date().getYear()
            y = parseInt("#20{vals[2]}")
            year = if cy>y then y else "#19{vals[2]}"
            "#{year}-#{vals[0]}-#{vals[1]}" 
        }
        {
          pattern: /^([1-9]|[12][0-9]|3[01])([- \/.])([1-9]|1[012])\2\d\d$/
          parse: (src, s) ->
            vals = src.split s
            cy = new Date().getYear()
            y = parseInt("#20{vals[2]}")
            year = if cy>y then y else "#19{vals[2]}"
            "#{year}-0#{vals[1]}-0#{vals[0]}" 
        }
        {
          pattern: /^\d\d([- \/.])(0[1-9]|1[012])\2(0[1-9]|[12][0-9]|3[01])$/
          parse: (src, s) ->
            vals = src.split s
            cy = new Date().getYear()
            y = parseInt("#20{vals[0]}")
            year = if cy>y then y else "#19{vals[0]}"
            "#{year}-#{vals[1]}-#{vals[2]}" 
        }
        {
          pattern: /^(0[1-9]|1[012])([- \/.])(0[1-9]|[12][0-9]|3[01])\2\d\d$/
          parse: (src, s) ->
            vals = src.split s
            cy = new Date().getYear()
            y = parseInt("#20{vals[2]}")
            year = if cy>y then y else "#19{vals[2]}"
            "#{year}-#{vals[0]}-#{vals[1]}" 
        }




        {
          pattern: /^(0[1-9]|[12][0-9]|3[01])([- \/.])(0[1-9]|1[012])\2\d\d$/
          parse: (src, s) ->
            vals = src.split s
            year = new Date().getYear()
            "#{year}-#{vals[1]}-#{vals[0]}" 
        }
        {
          pattern: /^\d\d([- \/.])(0[1-9]|1[012])\2(0[1-9]|[12][0-9]|3[01])$/
          parse: (src, s) ->
            vals = src.split s
            year = new Date().getYear()
            "#{year}-#{vals[1]}-#{vals[2]}" 
        }
        {
          pattern: /^(0[1-9]|1[012])([- \/.])(0[1-9]|[12][0-9]|3[01])\2\d\d$/
          parse: (src, s) ->
            vals = src.split s
            year = new Date().getYear()
            "#{year}-#{vals[0]}-#{vals[1]}" 
        }
      ]
    }
    {
      name: 'email'
      handlers: [
        {
          pattern: /\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}\b/i
        }
      ]
    }
  ]

  @getPattern: (src) ->
    for entity in @entities
      do (entity) ->
        for handler in entity.handlers
          do (handler) ->
            match = handler.pattern.exec src
            if match
              console.log match 
              console.log entity.name
              console.log handler.parse match

  @parse: (m) ->
    if m[3] 
        yyyy = m[2]
    else
        yyyy = if parseInt(m[2])<30 then "20#{m[2]}" else "19#{m[2]}"
    mm = if m[5] or not m[6] then m[4] else "0#{m[4]}"
    dd = if m[8] or not m[9] then m[7] else "0#{m[7]}"
    "#{yyyy}-#{mm}-#{dd}" 




