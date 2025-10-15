<template>
  <div
    class="dashboard"
    :class="{
      'sidebar-open': showInfoSidebar,
      'modal-open': forwardModalOpen || addMemberModalOpen,
    }"
  >
    <Sidebar
      :user="user"
      :conversation-updated="conversationUpdated"
      :active-conversation="activeConversation"
      @profile-update="$emit('profile-update', $event)"
      @active-conversation="onActiveConversation"
      @logout="$emit('logout')"
    />
    <Conversation
      :user="user"
      :conversation="activeConversation"
      :sidebar-open="showInfoSidebar"
      :forward-modal-open="forwardModalOpen"
      @message-sent="refreshConversations"
      @conversation-updated="refreshConversations"
      @forward-modal-open="onForwardModalOpen"
      @show-info-sidebar="handleShowInfoSidebar"
    />
    <ConversationInfoSidebar
      v-if="showInfoSidebar && activeConversation"
      :conversation="activeConversation"
      :user="user"
      @close="showInfoSidebar = false"
      @group-name-updated="updateGroupName"
      @group-updated="updateActiveConversation"
      @left-group="handleLeftGroup"
      @add-member-modal-open="addMemberModalOpen = true"
    />
  </div>
  <ForwardModal
    v-if="forwardModalOpen"
    :user="user"
    @close="forwardModalOpen = false"
    @forward="handleForward"
  />
  <AddMemberModal
    v-if="addMemberModalOpen"
    :conversation="activeConversation"
    @close="addMemberModalOpen = false"
    @member-added="updateActiveConversation"
  />
</template>

<script setup>
import { ref } from "vue";
import api from "@/services/api";

import Sidebar from "@/components/sidebar/Sidebar.vue";
import Conversation from "@/components/conversation/Conversation.vue";
import ConversationInfoSidebar from "@/components/conversation/ConversationInfoSidebar.vue";
import ForwardModal from "@/components/conversation/ForwardModal.vue";
import AddMemberModal from "@/components/conversation/AddMemberModal.vue";

const props = defineProps({
  user: Object,
});

const emit = defineEmits(["profile-update", "logout"]);

const activeConversation = ref(null);
const conversationUpdated = ref(Date.now());
const showInfoSidebar = ref(false);

function onActiveConversation(conversation) {
  showInfoSidebar.value = false;
  activeConversation.value = conversation;
  conversationUpdated.value = Date.now();
}

function refreshConversations() {
  conversationUpdated.value = Date.now();
}

function updateGroupName(newName) {
  activeConversation.value = {
    ...activeConversation.value,
    name: newName,
  };

  conversationUpdated.value = Date.now();
}

function updateActiveConversation(updatedConversation) {
  activeConversation.value = updatedConversation;
  conversationUpdated.value = Date.now();
}

function handleLeftGroup() {
  activeConversation.value = null;
  conversationUpdated.value = Date.now();
}

const forwardModalOpen = ref(false);
const messageToForward = ref(null);

function handleShowInfoSidebar() {
  if (!forwardModalOpen.value) showInfoSidebar.value = true;
}

function onForwardModalOpen({ open, message }) {
  forwardModalOpen.value = open;
  messageToForward.value = message;
}

async function handleForward(conversation) {
  if (!messageToForward.value) return;
  try {
    await api.post(`/conversations/${conversation.conversationId}/forwards`, {
      messageId: messageToForward.value.messageId,
    });

    forwardModalOpen.value = false;
    messageToForward.value = null;

    refreshConversations();
  } catch (e) {
    console.error(e);
  }
}

const addMemberModalOpen = ref(false);
</script>

<style scoped>
.dashboard {
  display: flex;
  width: 100vw;
  height: 100vh;
}

.dashboard.modal-open .conversation-info-sidebar,
.dashboard.modal-open .sidebar,
.dashboard.modal-open .conversation {
  opacity: 0.4;
  pointer-events: none;
  transition: opacity 0.2s;
}
</style>
