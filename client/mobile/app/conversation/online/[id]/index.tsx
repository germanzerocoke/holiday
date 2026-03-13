import { StyleSheet, Text, View } from "react-native";
import { useLocalSearchParams } from "expo-router";
import { colors } from "@/constants";
import { useEffect, useRef } from "react";
import { ConversationSignalResponse } from "@/types/conversation";
import {
  RTCPeerConnection,
  RTCIceCandidate,
  RTCSessionDescription,
} from "react-native-webrtc";

export default function OnlineConversationRoomScreen() {
  const { conversationId } = useLocalSearchParams();
  // const {
  //   data: conversation,
  //   isPending,
  //   isError,
  // } = useJoinOnlineConversation(String(roomId));

  const ws = useRef<WebSocket>(null);
  const peers = useRef<Record<string, RTCPeerConnection>>({});

  useEffect(() => {
    ws.current = new WebSocket(
      `ws://localhost:8080/online/conversation/join?conversationId=${conversationId}`,
    );

    ws.current.onmessage = async (event) => {
      const data: ConversationSignalResponse = JSON.parse(event.data);
      if (!data.signal) {
        for (const fromId of data.fromIds) {
          peers.current[fromId] = new RTCPeerConnection({
            iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
          });

          peers.current[fromId].addEventListener(
            "icecandidate",
            (event: RTCIceCandidate) => {
              if (event.candidate) {
                ws.current?.send(
                  JSON.stringify({
                    toId: fromId,
                    signal: { type: "candidate", candidate: event.candidate },
                  }),
                );
              }
            },
          );

          const offer = await peers.current[fromId].createOffer({
            offerToReceiveAudio: true,
            offerToReceiveVideo: false,
            voiceActivityDetection: true,
          });
          await peers.current[fromId].setLocalDescription(offer);
          ws.current?.send(
            JSON.stringify({
              toId: fromId,
              signal: peers.current[fromId].localDescription,
            }),
          );
        }
        return;
      }
      if (data.signal) {
        const fromId = data.fromIds[0];
        if (data.signal.type === "offer") {
          const offer = new RTCSessionDescription(data.signal);
          peers.current[fromId].setRemoteDescription(offer);
          const answer = await peers.current[fromId].createAnswer();
          peers.current[fromId].setLocalDescription(answer);
          ws.current?.send(
            JSON.stringify({
              toId: fromId,
              signal: peers.current[fromId].localDescription,
            }),
          );
          return;
        }
        if (data.signal.type === "answer") {
          const offerDescription = new RTCSessionDescription(data.signal);
          peers.current[fromId].setRemoteDescription(offerDescription);
          return;
        }
        if (data.signal.type === "candidate") {
          const iceCandidate = new RTCIceCandidate(data.signal.candidate);
          peers.current[fromId].addIceCandidate(iceCandidate);
          return;
        }
      }
    };

    return () => {
      ws.current?.close();
    };
  }, [conversationId]);

  return (
    <View style={styles.container}>
      <View style={styles.content}>
        <Text style={styles.title}></Text>
        <Text style={styles.description}></Text>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { backgroundColor: colors.SAND_110 },
  content: {
    padding: 16,
  },
  title: {
    fontSize: 18,
    color: colors.BLACK,
    fontWeight: 600,
    marginVertical: 8,
  },
  description: {
    fontSize: 16,
    color: colors.BLACK,
    marginBottom: 14,
  },
});
