define [], () ->
  class PGPatternMatcher
    ###
      Available Values Patterns:
      --------------------------
      Credit Card
      URL
      Email
      Domain
      IPv4
      UK National Insurance Number (NIN)
      UK Post Code
      VAT
      Currency (USA, EU and UK)
      Date (without time)
      Percent
      Float Number
      Integer Number
      Label (starting with a character and 32 max length)
      Text (any other stuff)

      Available Keys Patterns:
      --------------------------
      Identifier
      Date
      Coefficient
      Map Longitude
      Map Latitude

      THE ORDER MATTERS!! as the first matched pattern is returned

      TODO:
        more handlers for date entity
        include months names properly on date regex
        Datetime
        Address
        Bank Account
        Country codes
        ............
        more keys
    ###
    @entities: [

      # ------------------------------------ Credit Card Entity ----------------------------------
      {
        name: 'creditCard'
        handlers: [
          {
            pattern: ///
              ^
              (4\d{3}[ -]*\d{4}[ -]*\d{4}[ -]*\d(?:\d{3})?)|
              (5[1-5]\d{2}[ -]*\d{4}[ -]*\d{4}[ -]*\d{4})|
              (6(?:011|5[0-9]{2})[ -]*\d{4}[ -]*\d{4}[ -]*\d{4})|
              (3[47]\d{2}[ -]*\d{6}[ -]*\d{5})|
              (3(?:0[0-5]|[68][0-9])\d[ -]*\d{6}[ -]*\d{4})|
              ((?:2131|1800)[ -]*\d{6}[ -]*\d{5}|35\d{2}[ -]*\d{4}[ -]*\d{4}[ -]*\d{4})
              $
            ///i
          }
        ]
      }

      # ------------------------------------ Url Entity ----------------------------------
      {
        name: 'url'
        handlers: [
          {
            pattern: ///
              ^
              (# Scheme
               [a-z][a-z0-9+\-.]*:
               (# Authority & path
                //
                ([a-z0-9\-._~%!$&'()*+,;=]+@)?              # User
                ([a-z0-9\-._~%]+                            # Named host
                |\[[a-f0-9:.]+\]                            # IPv6 host
                |\[v[a-f0-9][a-z0-9\-._~%!$&'()*+,;=:]+\])  # IPvFuture host
                (:[0-9]+)?                                  # Port
                (/[a-z0-9\-._~%!$&'()*+,;=:@]+)*/?          # Path
               |# Path without authority
                (/?[a-z0-9\-._~%!$&'()*+,;=:@]+(/[a-z0-9\-._~%!$&'()*+,;=:@]+)*/?)?
               )
              |# Relative URL (no scheme or authority)
               ([a-z0-9\-._~%!$&'()*+,;=@]+(/[a-z0-9\-._~%!$&'()*+,;=:@]+)*/?  # Relative path
               |(/[a-z0-9\-._~%!$&'()*+,;=:@]+)+/?)                            # Absolute path
              )
              # Query
              (\?[a-z0-9\-._~%!$&'()*+,;=:@/?]*)?
              # Fragment
              (\#[a-z0-9\-._~%!$&'()*+,;=:@/?]*)?
              $
            ///i
          }
        ]
      }

      # ------------------------------------ Email Entity ----------------------------------
      {
        name: 'email'
        handlers: [
          {
            pattern: /[a-z0-9!#$%&'*+\/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?/i
          }
        ]
      }

      # ------------------------------------ Domain name Entity ----------------------------------
      {
        name: 'domain'
        handlers: [
          {
            pattern: /^\b((?=[a-z0-9-]{1,63}\.)(xn--)?[a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,63}\b$/i
          }
        ]
      }

      # ------------------------------------ IP v4 Entity ----------------------------------
      {
        name: 'ipv4'
        handlers: [
          {
            pattern: ///
              \b
              (25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.
              (25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.
              (25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.
              (25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])
              \b
            ///i
          }
        ]
      }

      # ------------------------------------ UK NIN Entity ----------------------------------
      {
        name: 'nin'
        handlers: [
          {
            pattern: /\b[abceghj-prstw-z][abceghj-nprstw-z] ?\d{2} ?\d{2} ?\d{2} ?[a-dfm]?\b/i
          }
        ]
      }

      # ------------------------------------ UK Post Code Entity ----------------------------------
      {
        name: 'postCode'
        handlers: [
          {
            pattern: /\b[A-Z]{1,2}[0-9][A-Z0-9]? [0-9][ABD-HJLNP-UW-Z]{2}\b/i
          }
        ]
      }

      # ------------------------------------ UK VAT Entity ----------------------------------
      {
        name: 'vat'
        handlers: [
          {
            pattern: /^(GB)?([0-9]{9}([0-9]{3})?|[A-Z]{2}[0-9]{3})$/i
          }
        ]
      }

      # ------------------------------------ EU/US Currency Entity ----------------------------------
      {
        name: 'currency'
        handlers: [
          {
            pattern: /^[+-]?[0-9]{1,3}(?:[0-9]*(?:[.,][0-9]{2})?|(?:,[0-9]{3})*(?:\.[0-9]{2})?|(?:\.[0-9]{3})*(?:,[0-9]{2})?)£€\$$/i
          }
        ]
      }

      # ------------------------------------ Date Entity ----------------------------------
      {
        name: 'date'
        handlers: [
          # ------------------------- Simple Numbered Date --------------------------
          {
            # Matches yyyy/mm/dd, yy/m/d, mm/dd .... with separators: - / . and space
            pattern: /^[^0-9]*(((19|20)?\d\d)[- \/])?((0)?([1-9])|1[012])[- \/]((0)?([1-9])|[12][0-9]|3[01])[^0-9]*$/
            # Groups: 1->yearEntry, 2->fullYear, 3->century, 4->monthEntry, 5->month0, 6->monthN, 7->dayEntry, 8->day0, 9->dayN
            parse: (m) => @parseDate m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9]
          }
          {
            # Matches yyyy/dd/mm, yy/d/m, dd/mm .... with separators: - / . and space
            pattern: /^[^0-9]*(((19|20)?\d\d)[- \/])?((0)?([1-9])|[12][0-9]|3[01])[- \/]((0)?([1-9])|1[012])[^0-9]*$/
            # Groups: 1->yearEntry, 2->fullYear, 3->century, 4->dayEntry, 5->day0, 6->dayN, 7->monthEntry, 8->month0, 9->monthN
            parse: (m) => @parseDate m[1], m[2], m[3], m[7], m[8], m[9], m[4], m[5], m[6]
          }
          {
            # Matches dd/mm/yyyy, d/m/yy, dd/mm .... with separators: - / . and space
            pattern: /^[^0-9]*((0)?([1-9])|[12][0-9]|3[01])[- \/]((0)?([1-9])|1[012])([- \/]((19|20)?\d\d))?[^0-9]*$/
            # Groups: 1->dayEntry, 2->day0, 3->dayN, 4->monthEntry, 5->month0, 6->monthN, 7->yearEntry, 8->fullYear, 9->century
            parse: (m) => @parseDate m[7], m[8], m[9], m[4], m[5], m[6], m[1], m[2], m[3]
          }
          {
            # Matches mm/dd/yyyy, m/d/yy, mm/dd .... with separators: - / . and space
            pattern: /^[^0-9]*((0)?([1-9])|1[012])[- \/]((0)?([1-9])|[12][0-9]|3[01])([- \/]((19|20)?\d\d))?[^0-9]*$/
            # Groups: 1->monthEntry, 2->month0, 3->monthN, 4->dayEntry, 5->day0, 6->dayN, 7->yearEntry, 8->fullYear, 9->century
            parse: (m) => @parseDate m[7], m[8], m[9], m[1], m[2], m[3], m[4], m[5], m[6]
          }
          # ------------------------- Named Months --------------------------
          {
            # Matches yyyy/mm/dd, yy/m/d, mm/dd .... with Month name and separators: - / . and space
            pattern: /^[^A-Z0-9]*(((19|20)?\d\d)[- \/]?)?([A-Z]{3,9})[- \/]?((0)?([1-9])|[12][0-9]|3[01])[^0-9]*$/i
            # Groups: 1->yearEntry, 2->fullYear, 3->century, 4->monthName, 5->dayEntry, 6->day0, 7->dayN
            parse: (m) =>
              mm = @parseMonthName m[4]
              return "error: month name" unless mm
              @parseDate m[1], m[2], m[3], mm, null, null, m[7], m[8], m[9]
          }
          {
            # Matches yyyy/dd/mm, yy/d/m, dd/mm .... with Month name and separators: - / . and space
            pattern: /^[^0-9]*(((19|20)?\d\d)[- \/])?((0)?([1-9])|[12][0-9]|3[01])[- \/]?([A-Z]{3,9})[^A-Z0-9]*$/i
            # Groups: 1->yearEntry, 2->fullYear, 3->century, 4->dayEntry, 5->day0, 6->dayN, 7->monthName
            parse: (m) =>
              mm = @parseMonthName m[7]
              return "error: month name" unless mm
              @parseDate m[1], m[2], m[3], mm, null, null, m[4], m[5], m[6]
          }
          {
            # Matches dd/mm/yyyy, d/m/yy, dd/mm .... with Month name and separators: - / . and space
            pattern: /^[^0-9]*((0)?([1-9])|[12][0-9]|3[01])[- \/]?([A-Z]{3,9})([- \/]?((19|20)?\d\d))?[^A-Z0-9]*$/i
            # Groups: 1->dayEntry, 2->day0, 3->dayN, 4->monthName, 5->yearEntry, 6->fullYear, 7->century
            parse: (m) =>
              mm = @parseMonthName m[4]
              return "error: month name" unless mm
              @parseDate m[7], m[8], m[9], mm, null, null, m[1], m[2], m[3]
          }
          {
            # Matches mm/dd/yyyy, m/d/yy, mm/dd .... with Month name and separators: - / . and space
            pattern: /^[^A-Z0-9]*([A-Z]{3,9})[- \/]?((0)?([1-9])|[12][0-9]|3[01])([- \/]((19|20)?\d\d))?[^0-9]*$/i
            # Groups: 1->monthName, 2->dayEntry, 3->day0, 4->dayN, 5->yearEntry, 6->fullYear, 7->century
            parse: (m) =>
              mm = @parseMonthName m[1]
              return "error: month name" unless mm
              @parseDate m[7], m[8], m[9], mm, null, null, m[4], m[5], m[6]
          }
        ]
      }

      # ------------------------------------ Integer Number Entity ----------------------------------
      {
        name: 'intNumber'
        handlers: [
          {
            pattern: /^[-+]?\d+$/i
            parse: (m) => parseInt m
          }
        ]
      }

      # ------------------------------------ Float Number Entity ----------------------------------
      {
        name: 'floatNumber'
        handlers: [
          {
            pattern: /^[-+]?\b[0-9]*\.?[0-9]+(?:[eE][-+]?[0-9]+)?\b$/i
            parse: (m) => parseFloat m
          }
        ]
      }

      # ------------------------------------ Percent Entity ----------------------------------
      {
        name: 'percent'
        handlers: [
          {
            pattern: /^\b[0-9]+(\.[0-9]+)?%\b$/i
            parse: (m) => parseFloat m
          }
        ]
      }

      # ------------------------------------ Label Entity ----------------------------------
      {
        name: 'label'
        handlers: [
          {
            pattern: /^[a-zA-Z]{1}.{0,63}$/i
          }
        ]
      }

      # ------------------------------------ Text Entity ----------------------------------
      {
        name: 'text'
        handlers: [
          {
            pattern: /^.*$/i
          }
        ]
      }
    ]

    @keys: [
      {
        name: 'identifier'
        pattern: /id|account|ref|((credit\b)?card)|(post\bcode)/i
      }
      {
        name: 'date'
        pattern: /date|year|day/i
      }
      {
        name: 'coefficient'
        pattern: /^coef|ind|ratio|percent|count$/i
      }
      {
        name: 'mapLongitude'
        pattern: /^(lon|ln|lng|long|longitud|longitude)$/i
      }
      {
        name: 'mapLatitude'
        pattern: /^(lat|lt|ltt|latit|latitud|latitude)$/i
      }
    ]

    # ------------------------------------ General functions ----------------------------------
    @getPattern: (src) ->
      pattern = null
      for entity in @entities
        do (entity) ->
          if not pattern
            for handler in entity.handlers
              do (handler) ->
                if not pattern
                  match = handler.pattern.exec src
                  if match
                    pattern = entity.name
                    console.log "Value Pattern --> #{pattern}: #{src}"
                    #console.log match
      pattern

    @getKeyPattern: (src) ->
      pattern = null
      for key in @keys
        do (key) ->
          if not pattern
            match = key.pattern.exec src
            if match
              pattern = key.name
              console.log "Key Pattern --> #{pattern}: #{src}"
              #console.log match
      pattern

    @parse: (src, pattern) ->
      parsed = null
      for entity in @entities when entity.name is pattern
        do (entity) ->
          for handler in entity.handlers
            do (handler) ->
              if not parsed
                match = handler.pattern.exec src
                if match
                  parsed = if handler.parse then handler.parse match else src
                #console.log match
                #console.log parsed
      parsed

    # ------------------------------------ Date functions ----------------------------------
    @parseMonthName: (name) ->
      mm = switch name
        when 'Jan', 'January', 'jan', 'january' then '01'
        when 'Feb', 'February', 'feb', 'february' then '02'
        when 'Mar', 'March', 'mar', 'march' then '03'
        when 'Apr', 'April', 'apr', 'april' then '04'
        when 'May', 'may' then '05'
        when 'Jun', 'June', 'jun', 'june' then '06'
        when 'Jul', 'July', 'jul', 'july' then '07'
        when 'Aug', 'August', 'aug', 'august' then '08'
        when 'Sep', 'September', 'sep', 'september' then '09'
        when 'Oct', 'October', 'oct', 'october' then '10'
        when 'Nov', 'November', 'nov', 'november' then '11'
        when 'Dec', 'December', 'dec', 'december' then '12'
        else null

    @parseDate: (ye, yyyy, yc, mm, m0, mn, dd, d0, dn) ->
      # If no year entry, assing current year
      yyyy = new Date().getYear() unless ye
      # Years with only 2 digits belongs to this century
      yyyy = "20#{yyyy}" if ye and not yc
      # Put '0' digit in month if it isn't already
      mm = "0#{mm}" if mn and not m0
      # Check days for months (no conditionals allowed in Javascript regex)
      switch parseInt(mm, 10)
        when 4,6,9,11
          return "error: month day" if parseInt(dd, 10)>30
        when 2
          nby = parseInt(yyyy, 10) % 4
          ddn = parseInt(dd, 10)
          return "error: month day" if ddn>29 or (ddn>28 and nby)
      # Put '0' digit in day if it isn't already
      dd = "0#{dd}" if dn and not d0
      # This will be our standard date format
      #"#{yyyy}-#{mm}-#{dd}"
      new Date(yyyy, mm-1, dd)
