<template>
  <div class="sidebar" @photo-updated="refreshUser">
    <SidebarHeader
      :user="user"
      :mode="mode"
      @show-conversation-list="mode = 'conversationList'"
      @show-new-conversation="showNewConversation"
      @show-profile="mode = 'profile'"
    />
    <ConversationList
      v-if="mode === 'conversationList'"
      :user="user"
      :conversation-updated="props.conversationUpdated"
      :active-conversation="props.activeConversation"
      @active-conversation="handleActiveConversation"
    />
    <NewConversation
      v-if="mode === 'newConversation'"
      @active-conversation="handleActiveConversation"
      @show-new-group="mode = 'newGroup'"
    />
    <NewGroup
      v-if="mode === 'newGroup'"
      @back="mode = 'newConversation'"
      @active-conversation="handleActiveConversation"
    />
    <Profile
      v-if="mode === 'profile'"
      :user="user"
      @profile-update="$emit('profile-update', $event)"
      @logout="$emit('logout')"
    />
  </div>
</template>

<script setup>
import { ref } from "vue";

import SidebarHeader from "@/components/sidebar/SidebarHeader.vue";
import ConversationList from "@/components/sidebar/ConversationList.vue";
import NewConversation from "@/components/sidebar/NewConversation.vue";
import NewGroup from "@/components/sidebar/NewGroup.vue";
import Profile from "@/components/sidebar/Profile.vue";

const props = defineProps({
  user: Object,
  activeConversation: Object,
  conversationUpdated: [Number, String, Boolean],
});

const emit = defineEmits(["profile-update", "active-conversation", "logout"]);

const mode = ref("conversationList");

function handleActiveConversation(conversation) {
  emit("active-conversation", conversation);
  mode.value = "conversationList";
}

function showNewConversation() {
  if (mode.value !== "newGroup") {
    mode.value = "newConversation";
  }
}
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--color-quaternary);
  width: 30vw;
  height: 100vh;
}
</style>
