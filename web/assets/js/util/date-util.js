const oneMinute = 1000 * 60; // 一分钟的毫秒数
const oneHour = oneMinute * 60; // 一小时的毫秒数
const oneDay = oneHour * 24; // 一天的毫秒数
const oneWeek = oneDay * 7; // 一星期的毫秒数
const oneMonth = oneDay * 30; // 一个月的毫秒数

/**
 * Decrease by days
 *
 * @param days the number of days to reduce
 */
Date.prototype.minusDays = function (days) {
    return this.minusMillis(oneDay * days);
};

/**
 * Increased by days
 *
 * @param days the number of days to add
 */
Date.prototype.plusDays = function (days) {
    return this.plusMillis(oneDay * days);
};

/**
 * Decrease by hour
 *
 * @param hours the number of hours to reduce
 */
Date.prototype.minusHours = function (hours) {
    return this.minusMillis(oneHour * hours);
};

/**
 * Increases by the hour
 *
 * @param hours the number of hours to add
 */
Date.prototype.plusHours = function (hours) {
    return this.plusMillis(oneHour * hours);
};

/**
 * Decrease by minute
 *
 * @param minutes the number of minutes to decrease
 */
Date.prototype.minusMinutes = function (minutes) {
    return this.minusMillis(oneMinute * minutes);
};

/**
 * Increments by minutes
 *
 * @param minutes the number of minutes to add
 */
Date.prototype.plusMinutes = function (minutes) {
    return this.plusMillis(oneMinute * minutes);
};

/**
 * decrements in milliseconds
 *
 * @param millis The number of milliseconds to reduce
 */
Date.prototype.minusMillis = function(millis) {
    let time = this.getTime() - millis;
    let newDate = new Date();
    newDate.setTime(time);
    return newDate;
};

/**
 * increments in milliseconds
 *
 * @param millis The number of milliseconds to add
 */
Date.prototype.plusMillis = function(millis) {
    let time = this.getTime() + millis;
    let newDate = new Date();
    newDate.setTime(time);
    return newDate;
};

/**
 * Set the time to 00:00:00.000 of the day
 */
Date.prototype.setMinTime = function () {
    this.setHours(0);
    this.setMinutes(0);
    this.setSeconds(0);
    this.setMilliseconds(0);
    return this;
};

/**
 * Set the time to 23:59:59.999 of the day
 */
Date.prototype.setMaxTime = function () {
    this.setHours(23);
    this.setMinutes(59);
    this.setSeconds(59);
    this.setMilliseconds(999);
    return this;
};

/**
 * format date
 */
Date.prototype.formatDate = function () {
    return this.getFullYear() + "-" + addZero(this.getMonth() + 1) + "-" + addZero(this.getDate());
};

/**
 * format time
 */
Date.prototype.formatTime = function () {
    return addZero(this.getHours()) + ":" + addZero(this.getMinutes()) + ":" + addZero(this.getSeconds());
};

/**
 * format date plus time
 *
 * @param split The separator between date and time, the default is a space
 */
Date.prototype.formatDateTime = function (split = ' ') {
    return this.formatDate() + split + this.formatTime();
};

class DateUtil {

    // 字符串转 Date 对象
    static parseDate(str) {
        return new Date(str.replace(/-/g, '/'));
    }

    static formatMillis(millis) {
        return moment(millis).format('YYYY-M-D H:m:s')
    }

    static firstDayOfMonth() {
        const date = new Date();
        date.setDate(1);
        date.setMinTime();
        return date;
    }
}