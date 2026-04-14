import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Button } from 'react-native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';

type RootStackParamList = {
  Home: undefined;
  Login: undefined;
  Register: undefined;
  NovelList: undefined;
  NovelDetail: { novelId: string };
  WriteChapter: { novelId: string };
  Chapter: { chapterId: string };
  Profile: undefined;
};

type ProfileScreenProps = {
  navigation: NativeStackNavigationProp<RootStackParamList, 'Profile'>;
};

// Mock user data for demonstration
const mockUser = {
  username: 'DemoUser',
  email: 'demo@example.com',
  joinedDate: '2024-01-01',
  novelsCreated: 3,
  chaptersWritten: 15,
};

export default function ProfileScreen({ navigation }: ProfileScreenProps) {
  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>{mockUser.username.charAt(0).toUpperCase()}</Text>
        </View>
        <Text style={styles.username}>{mockUser.username}</Text>
        <Text style={styles.email}>{mockUser.email}</Text>
      </View>
      
      <View style={styles.statsContainer}>
        <View style={styles.statItem}>
          <Text style={styles.statValue}>{mockUser.novelsCreated}</Text>
          <Text style={styles.statLabel}>Novels</Text>
        </View>
        <View style={styles.statItem}>
          <Text style={styles.statValue}>{mockUser.chaptersWritten}</Text>
          <Text style={styles.statLabel}>Chapters</Text>
        </View>
      </View>
      
      <View style={styles.infoContainer}>
        <Text style={styles.infoLabel}>Member since</Text>
        <Text style={styles.infoValue}>{mockUser.joinedDate}</Text>
      </View>
      
      <View style={styles.menuContainer}>
        <TouchableOpacity style={styles.menuItem}>
          <Text style={styles.menuText}>My Novels</Text>
        </TouchableOpacity>
        <TouchableOpacity style={styles.menuItem}>
          <Text style={styles.menuText}>Drafts</Text>
        </TouchableOpacity>
        <TouchableOpacity style={styles.menuItem}>
          <Text style={styles.menuText}>Settings</Text>
        </TouchableOpacity>
        <TouchableOpacity style={styles.menuItem}>
          <Text style={styles.menuText}>Help & Support</Text>
        </TouchableOpacity>
      </View>
      
      <View style={styles.buttonContainer}>
        <Button
          title="Logout"
          onPress={() => navigation.navigate('Home')}
        />
        <Button
          title="Back to Home"
          onPress={() => navigation.navigate('Home')}
        />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    padding: 20,
  },
  header: {
    alignItems: 'center',
    marginBottom: 30,
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: '#007AFF',
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 15,
  },
  avatarText: {
    fontSize: 32,
    color: '#fff',
    fontWeight: 'bold',
  },
  username: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  email: {
    fontSize: 14,
    color: '#666',
  },
  statsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 30,
    paddingVertical: 20,
    borderTopWidth: 1,
    borderBottomWidth: 1,
    borderColor: '#eee',
  },
  statItem: {
    alignItems: 'center',
  },
  statValue: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#007AFF',
  },
  statLabel: {
    fontSize: 14,
    color: '#666',
    marginTop: 5,
  },
  infoContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 20,
    paddingVertical: 10,
  },
  infoLabel: {
    fontSize: 14,
    color: '#666',
  },
  infoValue: {
    fontSize: 14,
    fontWeight: 'bold',
  },
  menuContainer: {
    marginBottom: 20,
  },
  menuItem: {
    paddingVertical: 15,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  menuText: {
    fontSize: 16,
  },
  buttonContainer: {
    gap: 10,
  },
});
