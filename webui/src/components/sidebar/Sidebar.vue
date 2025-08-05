<template>
    <div class="sidebar" @photo-updated="refreshUser">
        <SidebarHeader
            :user="user"
            :mode="mode"
            @showConversationList="mode = 'conversationList'"
            @showNewConversation="showNewConversation"
            @showProfile="mode = 'profile'"
        />
        <ConversationList
            v-if="mode === 'conversationList'"
            :user="user"
            :conversationUpdated="props.conversationUpdated"
            :activeConversation="props.activeConversation"
            @activeConversation="handleActiveConversation"
        />
        <NewConversation
            v-if="mode === 'newConversation'"
            @activeConversation="handleActiveConversation"
            @showNewGroup="mode = 'newGroup'"
        />
        <NewGroup
            v-if="mode === 'newGroup'"
            @back="mode = 'newConversation'"
            @activeConversation="handleActiveConversation"
        />
        <Profile
            :user="user"
            v-if="mode === 'profile'"
            @profile-update="$emit('profile-update', $event)"
        />
    </div>
</template>

<script setup>
import { ref } from "vue";

import SidebarHeader from "@/components/Sidebar/SidebarHeader.vue";
import ConversationList from "@/components/Sidebar/ConversationList.vue";
import NewConversation from "@/components/Sidebar/NewConversation.vue";
import NewGroup from "@/components/Sidebar/NewGroup.vue";
import Profile from "@/components/Sidebar/Profile.vue";

const props = defineProps({
    user: Object,
    activeConversation: Object,
    conversationUpdated: [Number, String, Boolean], // <-- aggiungi qui
});

const emit = defineEmits(["profile-update", "activeConversation"]);

const mode = ref("conversationList");

function handleActiveConversation(conversation) {
    emit("activeConversation", conversation);
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
