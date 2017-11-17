--liquibase formatted sql

--changeset id:1 author:mdyminski dbms:postgresql
CREATE TABLE str_users
(
  dispatchNumber   nvarchar(256) PRIMARY KEY,
  equipmentNumber  nvarchar(256),
  companyNumber    nvarchar(256),
  location         nvarchar(256),
  branch           nvarchar(256),
  round            nvarchar(256),
  visitTypeCode    nvarchar(256),
  customerNumber   nvarchar(256),
  status           nvarchar(256),
  eventCode        nvarchar(256),
  elevatorFloor    nvarchar(256),
  failureType      nvarchar(256),
  entrapment       BIT,
  personInjured    BIT,
  customerName     nvarchar(256),
  callerTel        nvarchar(256),
  lastModification datetimeoffset  NOT NULL
);
