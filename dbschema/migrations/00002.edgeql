CREATE MIGRATION m1h7phzyxodvonon4tb56k5bwjpjk6rxdwrsldvcta2dxbnu3dlixq
    ONTO m1ommioo7u3samgrrygru76rpdfuivif5u6yhq43rcfaychs5su63q
{
  CREATE TYPE default::HandInvoice {
      CREATE REQUIRED LINK message -> default::Message;
      CREATE REQUIRED PROPERTY comment -> std::str;
      CREATE REQUIRED PROPERTY total -> std::float64;
  };
};
