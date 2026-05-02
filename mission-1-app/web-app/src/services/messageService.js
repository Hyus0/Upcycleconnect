const MESSAGE_EVENT = "upcycle-message-change";

function currentUserId() {
  return localStorage.getItem("userId") || "guest";
}

function storageKey() {
  return `upcycle-messages:${currentUserId()}`;
}

function readConversations() {
  try {
    const parsed = JSON.parse(localStorage.getItem(storageKey()) || "[]");
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

function writeConversations(conversations) {
  localStorage.setItem(storageKey(), JSON.stringify(conversations));
  window.dispatchEvent(new Event(MESSAGE_EVENT));
}

export function getConversations() {
  return readConversations();
}

export function getConversation(conversationId) {
  return readConversations().find((conversation) => conversation.id === conversationId) || null;
}

export function startConversation(target) {
  const conversations = readConversations();
  const targetKey = `${target.kind || "contact"}:${target.targetId || target.name}:${target.contextId || "general"}`;
  const existing = conversations.find((conversation) => conversation.targetKey === targetKey);
  if (existing) return existing;

  const conversation = {
    id: `conv-${Date.now()}`,
    targetKey,
    kind: target.kind || "contact",
    targetId: target.targetId || null,
    name: target.name || "Contact",
    subject: target.subject || "Discussion",
    contextId: target.contextId || null,
    contextLabel: target.contextLabel || "",
    created_at: new Date().toISOString(),
    messages: [
      {
        id: `msg-${Date.now()}`,
        author: "system",
        body: `Conversation ouverte pour ${target.contextLabel || target.subject || "UpcycleConnect"}.`,
        created_at: new Date().toISOString()
      }
    ]
  };

  writeConversations([conversation, ...conversations]);
  return conversation;
}

export function sendMessage(conversationId, body) {
  const conversations = readConversations();
  const nextConversations = conversations.map((conversation) => {
    if (conversation.id !== conversationId) return conversation;
    return {
      ...conversation,
      messages: [
        ...conversation.messages,
        {
          id: `msg-${Date.now()}`,
          author: "me",
          body,
          created_at: new Date().toISOString()
        }
      ]
    };
  });
  writeConversations(nextConversations);
}

export function onMessagesChange(handler) {
  window.addEventListener(MESSAGE_EVENT, handler);
  window.addEventListener("storage", handler);
  return () => {
    window.removeEventListener(MESSAGE_EVENT, handler);
    window.removeEventListener("storage", handler);
  };
}
