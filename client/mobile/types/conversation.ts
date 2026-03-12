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

export interface ConversationFeedResponse {
  id: string;
  novel?: string;
  shortStory?: string;
  poem?: string;
  drama?: string;
  film?: string;
  by?: string;
  rule?: string;
  when: string;
  length: string;
  ongoing: boolean;
  isModerator: boolean;
  isRegistrant: boolean;
}

export interface OnlineConversationDetail {
  id: string;
  title: string;
  leaderId: string;
}
