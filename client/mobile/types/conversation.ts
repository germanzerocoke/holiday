export interface CreateOnlineConversationRequest {
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  by?: string;
  rule?: string;
  capacity: number;
  when: string;
  length: string;
}

export interface OnlineConversation {
  id: string;
  title: string;
}

export interface OnlineConversationDetail {
  id: string;
  title: string;
  leaderId: string;
}
