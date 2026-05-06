//// [tests/cases/compiler/declFileTypeofEnum.ts] ////

//// [declFileTypeofEnum.ts]
enum days {
    monday,
    tuesday,
    wednesday,
    thursday,
    friday,
    saturday,
    sunday
}

var weekendDay = days.saturday;
var daysOfMonth = days;
var daysOfYear: typeof days;

//// [declFileTypeofEnum.js]
"use strict";
var days;
(function (days) {
    days[days["monday"] = 0] = "monday";
    days[days["tuesday"] = 1] = "tuesday";
    days[days["wednesday"] = 2] = "wednesday";
    days[days["thursday"] = 3] = "thursday";
    days[days["friday"] = 4] = "friday";
    days[days["saturday"] = 5] = "saturday";
    days[days["sunday"] = 6] = "sunday";
})(days || (days = {}));
var weekendDay = days.saturday;
var daysOfMonth = days;
var daysOfYear;


//// [declFileTypeofEnum.d.ts]
enum days {
    monday = 0,
    tuesday = 1,
    wednesday = 2,
    thursday = 3,
    friday = 4,
    saturday = 5,
    sunday = 6
}
var weekendDay: days;
var daysOfMonth: typeof days;
var daysOfYear: typeof days;


//// [DtsFileErrors]


declFileTypeofEnum.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeofEnum.d.ts (1 errors) ====
    enum days {
    ~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        monday = 0,
        tuesday = 1,
        wednesday = 2,
        thursday = 3,
        friday = 4,
        saturday = 5,
        sunday = 6
    }
    var weekendDay: days;
    var daysOfMonth: typeof days;
    var daysOfYear: typeof days;
    