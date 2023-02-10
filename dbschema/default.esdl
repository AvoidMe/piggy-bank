module default {
  type User {
    required property username -> str;
    required property chat_id  -> int64 {
      constraint exclusive;
    };
  }

  type Message {
    required link     user -> User;
    required property text -> str;
    required property date -> datetime;
  }

  type Invoice {
    required link   message -> Message;
    required property total -> float64;
    required property  date -> datetime;
    
    required property   raw -> json;
  }

  type InvoiceItem {
    required property name     -> str;
    required property price    -> float64;
    required property quantity -> int64;

    required link     invoice  -> Invoice;
  }
};
