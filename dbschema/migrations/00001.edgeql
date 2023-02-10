CREATE MIGRATION m1ommioo7u3samgrrygru76rpdfuivif5u6yhq43rcfaychs5su63q
    ONTO initial
{
  CREATE TYPE default::User {
      CREATE REQUIRED PROPERTY chat_id -> std::int64 {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY username -> std::str;
  };
  CREATE TYPE default::Message {
      CREATE REQUIRED LINK user -> default::User;
      CREATE REQUIRED PROPERTY date -> std::datetime;
      CREATE REQUIRED PROPERTY text -> std::str;
  };
  CREATE TYPE default::Invoice {
      CREATE REQUIRED LINK message -> default::Message;
      CREATE REQUIRED PROPERTY date -> std::datetime;
      CREATE REQUIRED PROPERTY raw -> std::json;
      CREATE REQUIRED PROPERTY total -> std::float64;
  };
  CREATE TYPE default::InvoiceItem {
      CREATE REQUIRED LINK invoice -> default::Invoice;
      CREATE REQUIRED PROPERTY name -> std::str;
      CREATE REQUIRED PROPERTY price -> std::float64;
      CREATE REQUIRED PROPERTY quantity -> std::int64;
  };
};
