module.exports = (grunt) ->

  grunt.initConfig
    pkg: grunt.file.readJSON('package.json')
    destPath: 'public'
    srcPath: 'www-src'

    # CSS
    less:
      options:
        strictMath: false
        strictUnits: true
      css:
        options:
          sourceMap: true
          # sourceMapFilename: 'style.css.map'
        expand: true
        cwd: '<%= srcPath %>/less/'
        src: '*'
        dest: '<%= destPath %>/css/'
        ext: '.css'

    autoprefixer:
      css:
        options:
          map: true
        expand: true
        cwd: '<%= destPath %>/css/'
        src: '*'
        dest: '<%= destPath %>/css/'
        ext: '.css'
    # csslint:
    #     # The damn thing doesn't support the same config files as csslint itself
    #     # options:
    #     #     csslintrc: '.csslintrc'
    #   css:
    #     src: 'style.css'

    # JS
    coffeelint:
      options:
        configFile: "coffeelint.json"
      js:
        expand: true
        cwd: '<%= srcPath %>/coffee/'
        src: '*'

    coffee:
      options:
        sourceMap: true
      js:
        expand: true
        cwd: '<%= srcPath %>/coffee/'
        src: '*'
        dest: '<%= destPath %>/js/'
        ext: '.js'

    # Templates
    copy:
      templates:
        expand: true
        cwd: '<%= srcPath %>/templates/'
        src: '**/*'
        dest: '<%= destPath %>/'

    # Util
    clean:
      css:
        src: '<%= destPath %>/css'
      js:
        src: '<%= destPath %>/js'
      templates:
        src: ['<%= destPath %>/*.html', '<%=destPath %>/templates']

    watch:
      css:
        files: ['<%= srcPath %>/less/*']
        tasks: ['css']
      js:
        files: ['<%= srcPath %>/coffee/*']
        tasks: ['js']
      templates:
        files: ['<%= srcPath %>/templates/*']
        tasks: ['templates']

  # Builders
  grunt.loadNpmTasks 'grunt-contrib-less'
  grunt.loadNpmTasks 'grunt-contrib-coffee'
  grunt.loadNpmTasks 'grunt-autoprefixer'
  # grunt.loadNpmTasks 'grunt-hogan'

  # Complainers
  # grunt.loadNpmTasks 'grunt-contrib-csslint'
  grunt.loadNpmTasks 'grunt-coffeelint'

  # Util
  grunt.loadNpmTasks 'grunt-contrib-clean'
  grunt.loadNpmTasks 'grunt-contrib-watch'
  grunt.loadNpmTasks 'grunt-contrib-copy'
  # grunt.loadNpmTasks 'grunt-contrib-concat'
  # grunt.loadNpmTasks 'grunt-contrib-htmlmin'

  grunt.registerTask 'default', [
    'clean'
    'css'
    'js'
    'templates'
  ]

  grunt.registerTask 'css', [
    'less:css'
    'autoprefixer:css'
    # 'csslint:css'
  ]
  grunt.registerTask 'js', [
    # 'coffeelint:js'
    'coffee:js'
  ]
  grunt.registerTask 'templates', [
    # 'hogan:templates'
    'copy:templates'
  ]
