const gulp = require('gulp')
const bs = require('browser-sync').create()
const sass = require('gulp-sass')(require('sass'))
const run = require('gulp-run')

const path = require('path')

const dest = "./public/"
const src = "./static/"

const taskSass = (cb) => {
  return gulp.src(src + "scss/*.scss")
    .pipe(sass())
    .pipe(gulp.dest(dest + "css"))
    .pipe(bs.stream())
}

const taskGen = () => {
  return run('go run . pages').exec()
    .pipe(bs.stream())
}

const taskServe = gulp.series(taskSass, taskGen, function() {
  bs.init({
    server: dest
  })

  gulp.watch(path.join(dest, "!/public/**", "*.json"), taskGen)
  gulp.watch(path.join(src, "scss/*.scss"), taskSass)
  //gulp.watch(path.join(dest,"**/*.html")).on('change', bs.reload)
})

exports.default = taskServe
