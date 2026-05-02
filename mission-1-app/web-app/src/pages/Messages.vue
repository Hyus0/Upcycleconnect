<template>
  <main class="public-dashboard">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > MESSAGERIE</p>
        <h1 class="hero-title1">Messagerie</h1>
        <p class="classic-text">Echanges avec les vendeurs, formateurs et organisateurs.</p>
      </div>
    </header>

    <section class="messages-layout">
      <aside class="messages-list">
        <button
          v-for="conversation in conversations"
          :key="conversation.id"
          class="message-thread"
          :class="{ active: selectedId === conversation.id }"
          type="button"
          @click="selectedId = conversation.id"
        >
          <strong>{{ conversation.name }}</strong>
          <span>{{ conversation.contextLabel || conversation.subject }}</span>
        </button>

        <div v-if="conversations.length === 0" class="state-card">
          Aucune conversation. Utilisez le bouton contacter depuis une annonce, une formation ou un evenement.
        </div>
      </aside>

      <section class="messages-panel">
        <template v-if="selectedConversation">
          <div class="messages-panel__header">
            <div>
              <h2>{{ selectedConversation.name }}</h2>
              <p>{{ selectedConversation.contextLabel || selectedConversation.subject }}</p>
            </div>
            <span class="badge badge--green">{{ selectedConversation.kind }}</span>
          </div>

          <div class="messages-stream">
            <article
              v-for="message in selectedConversation.messages"
              :key="message.id"
              class="message-bubble"
              :class="{ 'message-bubble--me': message.author === 'me' }"
            >
              <p>{{ message.body }}</p>
              <small>{{ formatDate(message.created_at) }}</small>
            </article>
          </div>

          <form class="message-compose" @submit.prevent="handleSend">
            <textarea v-model="draft" rows="3" placeholder="Ecrire un message..." />
            <button class="btn-main-action" type="submit" :disabled="!draft.trim()">Envoyer</button>
          </form>
        </template>

        <div v-else class="state-card">Selectionnez une conversation.</div>
      </section>
    </section>
  </main>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import { getConversations, onMessagesChange, sendMessage } from "../services/messageService";

const route = useRoute();
const conversations = ref([]);
const selectedId = ref("");
const draft = ref("");
let stopSync = null;

const isLoggedIn = computed(() => !!localStorage.getItem("userToken"));
const userName = computed(() => {
  const prenom = localStorage.getItem("userPrenom") || "";
  const nom = localStorage.getItem("userNom") || "";
  return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});
const selectedConversation = computed(() => conversations.value.find((item) => item.id === selectedId.value));

function syncConversations() {
  conversations.value = getConversations();
  selectedId.value = route.query.conversation || selectedId.value || conversations.value[0]?.id || "";
}

function formatDate(value) {
  return new Date(value).toLocaleString("fr-FR", { day: "2-digit", month: "short", hour: "2-digit", minute: "2-digit" });
}

function handleSend() {
  if (!draft.value.trim() || !selectedId.value) return;
  sendMessage(selectedId.value, draft.value.trim());
  draft.value = "";
  syncConversations();
}

onMounted(() => {
  syncConversations();
  stopSync = onMessagesChange(syncConversations);
});

onBeforeUnmount(() => stopSync?.());
</script>

<style scoped>
.messages-layout {
  display: grid;
  grid-template-columns: minmax(260px, 360px) 1fr;
  gap: 18px;
}

.messages-list,
.messages-panel {
  min-height: 560px;
  border: 1px solid var(--border);
  border-radius: 16px;
  background: #fff;
  padding: 16px;
}

.messages-list {
  display: grid;
  align-content: start;
  gap: 10px;
}

.message-thread {
  display: grid;
  gap: 4px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: #f8fbf8;
  padding: 12px;
  text-align: left;
}

.message-thread.active {
  border-color: var(--brand-green);
  background: #edf8ef;
}

.message-thread span,
.messages-panel__header p,
.message-bubble small {
  color: var(--text-secondary);
}

.messages-panel {
  display: grid;
  grid-template-rows: auto 1fr auto;
  gap: 14px;
}

.messages-panel__header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  border-bottom: 1px solid var(--border);
  padding-bottom: 14px;
}

.messages-panel__header h2,
.messages-panel__header p {
  margin: 0;
}

.messages-stream {
  display: grid;
  align-content: start;
  gap: 10px;
  overflow: auto;
}

.message-bubble {
  width: min(78%, 620px);
  border-radius: 14px;
  background: #f2f6f3;
  padding: 12px 14px;
}

.message-bubble--me {
  justify-self: end;
  background: #dff2e5;
}

.message-bubble p {
  margin: 0 0 6px;
}

.message-compose {
  display: flex;
  gap: 12px;
}

.message-compose textarea {
  flex: 1;
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 12px;
  resize: vertical;
}

@media (max-width: 860px) {
  .messages-layout {
    grid-template-columns: 1fr;
  }
}
</style>
