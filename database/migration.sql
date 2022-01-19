CREATE TABLE `module`
(
    moduleCode varchar(255),
    moudleName varchar(255),
    moduleSynopsis varchar(255),
    moduleLO varchar(255),
    moudleClasses varchar(255),
    tutor varchar(255),
    students varchar(255),
    moudleRnC varchar(255),
    primary key (moduleCode)
);

INSERT INTO `module` (`moduleCode`,`moudleName`,`moduleSynopsis`,`moduleLO`,`moudleClasses`,`tutor`,`students`,`moudleRnC`)
VALUE ('ETI01','ETI','learn about ERI', 'prepare for new trent', 'TO1', 'Wesely Teo','xiong run lin', 'NA');
