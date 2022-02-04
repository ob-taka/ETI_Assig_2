-- invoke table
CREATE TABLE `module`
(
    moduleCode varchar(255),
    moudleName varchar(255),
    moduleSynopsis varchar(255),
    moduleLO varchar(255),
    moudleRnC varchar(255),
    primary key (moduleCode)
);

CREATE TABLE `tutor`
(
    tID varchar(255),
    tMoulde varchar(255),
    primary key (tID)
);

CREATE TABLE `student`
(
    stuID varchar(255),
    sModuel varchar(255),
    primary key (stuID)
);

CREATE TABLE `class`
(
    classCode varchar(255),
    mouduleCode varchar(255),
    primary key (classCode)
);

-- Insert vaule into table
INSERT INTO `module` (`moduleCode`,`moudleName`,`moduleSynopsis`,`moduleLO`,`moudleRnC`)
VALUE ('ETI01','ETI','learn about ERI', 'prepare for new trent', 'NA');
INSERT INTO `module` (`moduleCode`,`moudleName`,`moduleSynopsis`,`moduleLO`,`moudleRnC`)
VALUE ('DF01','DF','learn about DF','Digital investigation','NA');

INSERT INTO `tutor` (`tID`,`tMoulde`)
VALUE ('T123', 'ETI01');
INSERT INTO `tutor` (`tID`,`tMoulde`)
VALUE ('T223', 'ETI01');

INSERT INTO `student` (`stuID`,`sModuel`)
VALUE ('S123', 'ETI01');
INSERT INTO `student` (`stuID`,`sModuel`)
VALUE ('S223', 'ETI01');

INSERT INTO `class` (`classCode`,`mouduleCode`)
VALUE ('C01', 'ETI01');
INSERT INTO `class` (`classCode`,`mouduleCode`)
VALUE ('C02', 'ETI01');



