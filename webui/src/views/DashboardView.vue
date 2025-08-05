<template>
    <div class="dashboard">
        <Sidebar
            :user="user"
            :conversationUpdated="conversationUpdated"
            :activeConversation="activeConversation"
            @profile-update="$emit('profile-update', $event)"
            @activeConversation="onActiveConversation"
        />
        <Conversation
            :user="user"
            :conversation="activeConversation"
            @messageSent="refreshConversations"
            @conversationUpdated="refreshConversations"
        />
    </div>
</template>

<script setup>
import { ref } from "vue";

import Sidebar from "@/components/sidebar/Sidebar.vue";
import Conversation from "@/components/conversation/Conversation.vue";

const props = defineProps({
    user: Object,
});

const emit = defineEmits(["profile-update"]);

const activeConversation = ref(null);
const conversationUpdated = ref(Date.now());

function onActiveConversation(conversation) {
    activeConversation.value = conversation;
    conversationUpdated.value = Date.now();
}

function refreshConversations() {
    conversationUpdated.value = Date.now();
}
</script>

<style scoped>
.dashboard {
    display: flex;
    height: 100vh;
}
</style>
